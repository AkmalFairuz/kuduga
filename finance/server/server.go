package server

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/digiflazz"
	"github.com/akmalfairuz/finance/module/duitku"
	"github.com/akmalfairuz/finance/server/handler"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/migration"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/akmalfairuz/finance/server/service/paymentmethod"
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	config Config

	start sync.Once

	log         *logrus.Logger
	db          *database.DB
	app         *fiber.App
	firebaseApp *firebase.App
	redisDb     *redis.Client
	scheduler   gocron.Scheduler
}

func New(log *logrus.Logger, config Config) *Server {
	log.Infof("Connecting to the mysql database...")
	db, err := database.New(config.Database.DSN(), "mysql")
	if err != nil {
		log.Fatalf("error connecting to mysql database: %v", err)
	}

	if err := migration.Up(log, db); err != nil {
		log.Fatal(err)
	}

	log.Infof("Connecting to the redis database...")
	redisDb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
	})
	redisStatus := redisDb.Ping(context.TODO())
	if redisStatus.Err() != nil {
		log.Fatalf("error connecting to redis database: %+v", redisStatus.Err())
	}

	redisLocker, err := redislock.NewRedisLocker(redisDb)
	if err != nil {
		log.Fatalf("error creating redis locker for scheduler: %v", err)
	}
	scheduler, err := gocron.NewScheduler(gocron.WithDistributedLocker(redisLocker))
	if err != nil {
		log.Fatalf("error creating scheduler: %v", err)
	}

	log.Infof("Connecting to firebase service...")
	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to connect to firebase service: %+v", err)
	}

	log.Infof("Connecting to firebase cloud messaging...")
	firebaseMessaging, err := firebaseApp.Messaging(context.TODO())
	if err != nil {
		log.Fatalf("failed to connect to firebase cloud messaing: %+v", err)
	}

	minioStorage, err := provider.NewMinioStorageProvider(config.Minio)
	if err != nil {
		log.Fatalf("failed to initialize minio client: %v", err)
	}

	fiberConfig := fiber.Config{
		AppName:       "finance",
		CaseSensitive: true,
		ErrorHandler:  middleware.ErrorHandlerMiddleware(),
		ReadTimeout:   time.Second * 30,
		WriteTimeout:  time.Second * 30,
		IdleTimeout:   time.Minute * 15,
	}

	if config.IsProduction() {
		fiberConfig.ProxyHeader = "X-Real-IP"
	}

	app := fiber.New(fiberConfig)

	dgClient := digiflazz.New(config.Digiflazz)

	userRepository := repository.NewUserRepository(db)
	otpRepository := repository.NewOTPRepository(db)
	tokenRepository := repository.NewTokenRepository(db)
	resetPasswordRepository := repository.NewResetPasswordRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	depositRepository := repository.NewDepositRepository(db)
	productRepository := repository.NewProductRepository(db)
	purchaseRepository := repository.NewPurchaseRepository(db)
	metaRepository := repository.NewMetaRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)
	transferRepository := repository.NewTransferRepository(db)
	supportTicketRepository := repository.NewSupportTicketRepository(db)
	checkerRepository := repository.NewCheckerRepository(redisDb)
	userKycRepository := repository.NewUserKycRepository(db)

	authService := service.NewAuthService(log.WithField("service", "auth"), tokenRepository)
	userService := service.NewUserService(log, userRepository, resetPasswordRepository)
	emailService := service.NewEmailService(provider.NewAmazonSESEMailProvider(config.AWSCredentials))
	otpService := service.NewOTPService(log.WithField("service", "otp"), config.OTPEmailFrom, redisDb, otpRepository)
	paymentService := service.NewPaymentService(log.WithField("service", "payment"), scheduler, depositRepository)
	transactionService := service.NewTransactionService(log.WithField("service", "transaction"), transactionRepository)
	productService := service.NewProductService(productRepository)
	purchaseProcessorService := service.NewPurchaseProcessorService(log.WithField("service", "purchaseProcessor"), scheduler)
	purchaseService := service.NewPurchaseService(log, purchaseRepository)
	notificationService := service.NewNotificationService(notificationRepository, provider.NewFirebaseNotificationProvider(firebaseMessaging))
	checkerService := service.NewCheckerService(checkerRepository)
	metaService := service.NewMetaService(metaRepository)
	transferService := service.NewTransferService(log.WithField("service", "transfer"), transferRepository)
	supportService := service.NewSupportService(supportTicketRepository, minioStorage, redisDb)
	productExternalManagementService := service.NewProductExternalManagementService(log.WithField("service", "productExternalManagement"), scheduler, dgClient)
	userKycService := service.NewUserKycService(userKycRepository, minioStorage)
	analyticsService := service.NewAnalyticsService()

	app.Use(middleware.RequestLogMiddleware(log))

	// use network throttling
	//app.Use(func(c *fiber.Ctx) error {
	//	time.Sleep(500 * time.Millisecond)
	//	return c.Next()
	//})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORSAllowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	authService.SetGoogleOauthClientId(config.GoogleOAuthClientId)

	otpService.SetEmailService(emailService)

	paymentService.SetTransactionService(transactionService)
	paymentService.SetNotificationService(notificationService)
	if config.DiscordWebhook.Deposit != "" {
		paymentService.SetDepositAlert(provider.NewDiscordTextAlertProvider(config.DiscordWebhook.Deposit))
	}

	transactionService.SetUserService(userService)

	notificationService.SetAuthService(authService)

	purchaseProcessorService.RegisterProvider(provider.DigiflazzPurchaseProviderName, provider.NewDigiflazzPurchaseProvider(log, dgClient))

	purchaseService.SetProductService(productService)
	purchaseService.SetPaymentService(paymentService)
	purchaseService.SetTransactionService(transactionService)
	purchaseService.SetPurchaseProcessorService(purchaseProcessorService)
	purchaseService.SetNotificationService(notificationService)
	purchaseService.SetCheckerService(checkerService)
	if config.DiscordWebhook.Purchase != "" {
		purchaseService.SetAlert(provider.NewDiscordTextAlertProvider(config.DiscordWebhook.Purchase))
	}

	purchaseProcessorService.SetPurchaseService(purchaseService)

	transferService.SetUserService(userService)
	transferService.SetTransactionService(transactionService)
	transferService.SetNotificationService(notificationService)

	productExternalManagementService.SetProductService(productService)
	if config.DiscordWebhook.ProductManagement != "" {
		productExternalManagementService.SetAlert(provider.NewDiscordTextAlertProvider(config.DiscordWebhook.ProductManagement))
	}

	generalAlert := provider.NewDiscordTextAlertProvider(config.DiscordWebhook.General)

	userKycService.SetNotificationService(notificationService)
	userKycService.SetUserService(userService)
	userKycService.SetAlert(generalAlert)

	supportService.SetAlert(generalAlert)
	userService.SetAlert(generalAlert)

	authHandler := handler.NewAuthHandler(authService, userService)
	authHandler.Route(app)

	userHandler := handler.NewUserHandler(userService, authService, otpService, userKycService)
	userHandler.Route(app)

	balanceHandler := handler.NewBalanceHandler(authService, paymentService, transactionService, userService, transferService)
	balanceHandler.Route(app)

	callbackHandler := handler.NewCallbackHandler(purchaseProcessorService, paymentService)
	callbackHandler.SetBcaMutasiApiKey(config.BcaPayment.CallbackApiKey)
	callbackHandler.Route(app)

	purchaseHandler := handler.NewPurchaseHandler(purchaseService, authService, userService)
	purchaseHandler.Route(app)

	productHandler := handler.NewProductHandler(productService, metaService)
	productHandler.Route(app)

	adminHandler := handler.NewAdminHandler(log, authService, productService, notificationService, purchaseService, metaService, productExternalManagementService, supportService, userKycService, userService, purchaseProcessorService, analyticsService)
	adminHandler.Route(app)

	notificationHandler := handler.NewNotificationHandler(authService, notificationService)
	notificationHandler.Route(app)

	checkerHandler := handler.NewCheckerHandler(checkerService, authService)
	checkerHandler.Route(app)

	appHandler := handler.NewAppHandler(metaService)
	appHandler.Route(app)

	supportHandler := handler.NewSupportHandler(supportService, authService, userService)
	supportHandler.Route(app)

	if config.AppMode == "development" {
		log.Infof("App running in development mode")
		testingHandler := handler.NewTestingHandler(otpService)
		testingHandler.Route(app)
	}

	duitkuClient := duitku.NewClient(log, config.Duitku)
	callbackHandler.SetDuitkuClient(duitkuClient)
	// Transfer Bank
	paymentService.RegisterPaymentMethod(paymentmethod.NewBankBCA(config.BcaPayment))
	// QR Codes
	paymentService.RegisterPaymentMethod(paymentmethod.NewQris(duitkuClient))
	// Virtual Account
	paymentService.RegisterPaymentMethod(paymentmethod.NewBRIVA(duitkuClient))
	paymentService.RegisterPaymentMethod(paymentmethod.NewPermataBankVirtual(duitkuClient))
	paymentService.RegisterPaymentMethod(paymentmethod.NewBNIVirtual(duitkuClient))
	//paymentService.RegisterPaymentMethod(paymentmethod.NewATMBersama(duitkuClient))
	// Retail
	paymentService.RegisterPaymentMethod(paymentmethod.NewAlfamart(duitkuClient))

	analyticsService.SetPaymentService(paymentService)
	analyticsService.SetAuthService(authService)
	analyticsService.SetPurchaseService(purchaseService)
	analyticsService.SetUserService(userService)

	userService.SetEmailService(emailService)

	return &Server{
		config:      config,
		log:         log,
		db:          db,
		app:         app,
		firebaseApp: firebaseApp,
		redisDb:     redisDb,
		scheduler:   scheduler,
	}
}

func (s *Server) Start() {
	s.start.Do(func() {
		s.log.Infof("Starting scheduler...")
		s.scheduler.Start()

		s.log.Infof("Starting server...")

		s.app.Use(func(c *fiber.Ctx) error {
			return c.Status(http.StatusNotFound).JSON(responses.M("Not found"))
		})

		if err := s.app.Listen(s.config.ListenAddress); err != nil {
			s.log.Fatalf("Failed to listen on %v: %v", s.config.ListenAddress, err)
		}
	})
}

func (s *Server) CloseOnProgramEnd() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		s.Close()
	}()
}

func (s *Server) Close() {
	s.log.Infof("Closing server...")
	if err := s.app.Shutdown(); err != nil {
		s.log.Errorf("Failed to shutdown server: %v", err)
	}
	if err := s.scheduler.Shutdown(); err != nil {
		s.log.Errorf("Failed to shutdown scheduler: %v", err)
	}
	if err := s.db.Close(); err != nil {
		s.log.Errorf("Failed to close database connection: %v", err)
	}
	if err := s.redisDb.Close(); err != nil {
		s.log.Errorf("Failed to close redis connection: %v", err)
	}
	s.log.Infof("Server closed")
}
