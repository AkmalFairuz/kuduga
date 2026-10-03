package handler

import (
	"errors"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type UserHandler struct {
	userService    *service.UserService
	userKycService *service.UserKycService
	authService    *service.AuthService
	otpService     *service.OTPService
}

func NewUserHandler(userService *service.UserService, authService *service.AuthService, otpService *service.OTPService, userKycService *service.UserKycService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		authService:    authService,
		otpService:     otpService,
		userKycService: userKycService,
	}
}

func (h *UserHandler) Route(r fiber.Router) {
	group := r.Group("/user")

	group.Post("/register", h.handleRegister)
	group.Post("/registerOTP", h.handleRegisterOTP)
	group.Post("/registerWithGoogle", h.handleRegisterWithGoogle)

	group.Post("/resetPasswordOTP", h.handleResetPasswordOTP)
	group.Post("/resetPasswordVerifyOTP", h.handleResetPasswordVerifyOTP)
	group.Post("/resetPassword", h.handleResetPassword)

	group.Post("/resendOTP", h.handleResendOTP)

	auth := middleware.AuthMiddleware(h.authService)

	group.Post("/updateDeviceUniqueID", auth, h.handleUpdateDeviceUniqueID)
	group.Post("/updateEmailOTP", auth, h.handleUpdateEmailOTP)
	group.Post("/updateEmail", auth, h.handleUpdateEmail)
	group.Post("/updateDisplayName", auth, h.handleUpdateDisplayName)
	group.Post("/updatePassword", auth, h.handleUpdatePassword)

	pinGroup := group.Group("/pin", auth)
	pinGroup.Post("/create", h.handleCreatePin)
	pinGroup.Post("/update", h.handleUpdatePin)
	pinGroup.Post("/delete", h.handleDeletePin)
	pinGroup.Post("/recoveryOTP", h.handleRecoveryOTPPin)
	pinGroup.Post("/recovery", h.handleRecoveryPin)

	kycGroup := group.Group("/kyc", auth)
	kycGroup.Get("/status", h.handleKycStatus)
	kycGroup.Post("/request", h.handleKycRequest)

	referralGroup := group.Group("/referral", auth)
	referralGroup.Get("/info", h.handleReferralInfo)
	referralGroup.Post("/use", h.handleUseReferral)
}

func (h *UserHandler) handleUpdateDeviceUniqueID(ctx *fiber.Ctx) error {
	var req requesttype.UpdateUserDeviceUniqueIDRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)
	if err := h.userService.UpdateUserDeviceUniqueID(userId, req.DeviceUniqueID); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Device unique ID updated"))
}

func (h *UserHandler) handleUseReferral(ctx *fiber.Ctx) error {
	var req requesttype.UseReferralCodeRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}
	if user.ReferralUserID != nil {
		return errortype.ErrReferralAlreadyUsed
	}
	if err := h.userService.UseReferralCode(&user, req.ReferralCode); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Success using referral code"))
}

func (h *UserHandler) handleReferralInfo(ctx *fiber.Ctx) error {
	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}
	if user.ReferralCode == nil {
		if _, err := h.userService.GenerateReferralCode(&user); err != nil {
			return err
		}
	}
	totalUserUsingYourReferralCode, err := h.userService.GetCountUsedThisReferralUser(user.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(&responsetype.UserReferralInfoResponse{
		Message:                        "Undang teman anda untuk menggunakan aplikasi Kuduga!",
		YourReferralCode:               *user.ReferralCode,
		TotalUserUsingYourReferralCode: totalUserUsingYourReferralCode,
		ShareText:                      "Download aplikasi Kuduga untuk melakukan pembelian pulsa, paket data dan topup game dengan harga yang murah di https://example.invalid/dl?ref=" + *user.ReferralCode + " dan gunakan kode referral " + *user.ReferralCode + " saat mendaftar.",
	})
}

func (h *UserHandler) handleKycStatus(ctx *fiber.Ctx) error {
	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}
	kyc, err := h.userKycService.GetByUserID(user.ID)
	var requestStatus *int
	if err != nil && !errors.Is(err, database.ErrNoRows) {
		return err
	}
	if err == nil {
		requestStatus = &kyc.Status
	}
	return ctx.JSON(&responsetype.UserKycStatusResponse{
		CurrentStatus: user.KycStatus,
		RequestStatus: requestStatus,
	})
}

func (h *UserHandler) handleKycRequest(ctx *fiber.Ctx) error {
	var req requesttype.KycRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	files, err := ctxhelper.ReadFiles(ctx)
	if err != nil {
		return err
	}

	var documentFile *model.File
	for _, file := range files {
		if file.Name == "ktp.png" {
			documentFile = file
			break
		}
	}

	if documentFile == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.M("Document required"))
	}

	create := &model.CreateUserKyc{
		UserID:       ctxhelper.GetUserID(ctx),
		FullName:     req.FullName,
		DocumentID:   req.DocumentID,
		DocumentFile: documentFile.Bytes,
	}

	if err := h.userKycService.Create(create); err != nil {
		return err
	}

	return ctx.JSON(responses.M("KYC request sent"))
}

func (h *UserHandler) handleRecoveryOTPPin(ctx *fiber.Ctx) error {
	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}
	rawTok := random.Token()
	otp, err := h.otpService.RequestEmailOTP(ctxhelper.SafeUniqueAddress(ctx), user.Email, model.OTPScopeRecoveryPin, rawTok)
	if err != nil {
		return err
	}
	return ctx.JSON(responsetype.NewOTPInfoFromModel(otp, rawTok))
}

func (h *UserHandler) handleRecoveryPin(ctx *fiber.Ctx) error {
	var req requesttype.OTPInfo
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}
	if err := h.otpService.VerifyOTP(req.ID, user.Email, model.OTPScopeRecoveryPin, req.Token, req.Code); err != nil {
		return err
	}
	if err := h.userService.DeleteUserPin(user.ID, "", true); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Pin deleted"))
}

func (h *UserHandler) handleIsPinEnabled(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)
	user, err := h.userService.GetUser(userId)
	if err != nil {
		return err
	}
	return ctx.JSON(&responsetype.IsPinEnabledResponse{Enabled: user.PinHash != nil})
}

func (h *UserHandler) handleDeletePin(ctx *fiber.Ctx) error {
	var req requesttype.DeletePinRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)
	if err := h.userService.DeleteUserPin(userId, req.Pin, false); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Pin deleted"))
}

func (h *UserHandler) handleCreatePin(ctx *fiber.Ctx) error {
	var req requesttype.CreatePinRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)
	if err := h.userService.CreateUserPin(userId, req.Pin); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Pin created"))
}

func (h *UserHandler) handleUpdatePin(ctx *fiber.Ctx) error {
	var req requesttype.UpdatePinRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)
	if err := h.userService.UpdateUserPin(userId, req.Pin, req.NewPin); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Pin updated"))
}

func (h *UserHandler) handleUpdatePassword(ctx *fiber.Ctx) error {
	var req requesttype.UpdateUserPasswordRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	tok := ctxhelper.GetToken(ctx)
	userId := tok.UserID

	if err := h.userService.VerifyPassword(userId, req.CurrentPassword); err != nil {
		return fiber.NewError(http.StatusUnauthorized, "wrong password")
	}

	if err := h.userService.UpdateUserPassword(userId, req.NewPassword); err != nil {
		return err
	}

	tokens, err := h.authService.GetTokens(userId)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if token.ID == tok.ID {
			continue
		}
		if err := h.authService.RevokeToken(token.ID); err != nil {
			return err
		}
	}

	return ctx.JSON(responsetype.MessageResponse{
		Message: "Password updated",
	})
}

func (h *UserHandler) handleUpdateEmailOTP(ctx *fiber.Ctx) error {
	var req requesttype.UpdateUserEmailOTPRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.userService.VerifyPassword(ctxhelper.GetUserID(ctx), req.Password); err != nil {
		return err
	}

	rawTok := random.Token()
	otp, err := h.otpService.RequestEmailOTP(ctxhelper.SafeUniqueAddress(ctx), req.Email, model.OTPScopeUpdateEmail, rawTok)
	if err != nil {
		return err
	}

	return ctx.JSON(responsetype.NewOTPInfoFromModel(otp, rawTok))
}

func (h *UserHandler) handleUpdateDisplayName(ctx *fiber.Ctx) error {
	var req requesttype.UpdateUserDisplayNameRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	userId := ctxhelper.GetUserID(ctx)

	if err := h.userService.UpdateUserDisplayName(userId, req.DisplayName); err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) handleUpdateEmail(ctx *fiber.Ctx) error {
	var req requesttype.UpdateUserEmailRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.otpService.VerifyOTP(req.OTP.ID, req.Email, model.OTPScopeUpdateEmail, req.OTP.Token, req.OTP.Code); err != nil {
		return err
	}

	userId := ctxhelper.GetUserID(ctx)

	if err := h.userService.UpdateUserEmail(userId, req.Email); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Email updated"))
}

func (h *UserHandler) handleResendOTP(ctx *fiber.Ctx) error {
	var req requesttype.ResendOTPRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.otpService.ResendOTP(req.ID, req.Token); err != nil {
		return err
	}

	return ctx.JSON(responses.M("OTP resent"))
}

func (h *UserHandler) handleResetPasswordOTP(ctx *fiber.Ctx) error {
	var req requesttype.ResetPasswordOTPRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if _, err := h.userService.GetUserByEmail(req.Email); err != nil {
		return err
	}

	rawTok := random.Token()
	otp, err := h.otpService.RequestEmailOTP(ctxhelper.SafeUniqueAddress(ctx), req.Email, model.OTPScopeResetPassword, rawTok)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(&responsetype.ResetPasswordOTPResponse{OTP: responsetype.NewOTPInfoFromModel(otp, rawTok)})
}

func (h *UserHandler) handleResetPasswordVerifyOTP(ctx *fiber.Ctx) error {
	var req requesttype.ResetPasswordVerifyOTPRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.otpService.VerifyOTP(req.OTP.ID, req.Email, model.OTPScopeResetPassword, req.OTP.Token, req.OTP.Code); err != nil {
		return err
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	token := random.Token()
	resetPassword, err := h.userService.CreateResetPassword(user.ID, token)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(&responsetype.ResetPasswordVerifyOTPResponse{ResetPasswordID: resetPassword.ID, Token: token})
}

func (h *UserHandler) handleResetPassword(ctx *fiber.Ctx) error {
	var req requesttype.ResetPasswordRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.userService.ResetPassword(req.ResetPasswordID, req.Token, req.NewPassword); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Password changed"))
}

func (h *UserHandler) handleRegisterOTP(ctx *fiber.Ctx) error {
	var req requesttype.RegisterOTPRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	{
		if !model.IsValidUsername(req.Name) {
			return ctx.Status(http.StatusBadRequest).JSON(responsetype.RegisterFail{Message: "Username hanya boleh mengandung huruf, angka, dan underscore. Panjang username minimal 3 karakter dan maksimal 20 karakter"})
		}
	}
	{
		_, err := h.userService.GetUserByEmail(req.Email)
		if err == nil {
			return ctx.Status(http.StatusConflict).JSON(responsetype.RegisterFail{Message: "Alamat email '" + req.Email + "' sudah digunakan di akun lain. Mohon gunakan email yang belum terdaftar"})
		}
		if !errors.Is(err, database.ErrNoRows) {
			return err
		}
	}
	{
		_, err := h.userService.GetUserByName(req.Name)
		if err == nil {
			return ctx.Status(http.StatusConflict).JSON(responsetype.RegisterFail{Message: "Username '" + req.Name + "' tidak tersedia. Mohon untuk mengganti username yang belum tersedia"})
		}
		if !errors.Is(err, database.ErrNoRows) {
			return err
		}
	}

	rawTok := random.Token()
	otp, err := h.otpService.RequestEmailOTP(ctxhelper.SafeUniqueAddress(ctx), req.Email, model.OTPScopeRegister, rawTok)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(responsetype.NewOTPInfoFromModel(otp, rawTok))
}

func (h *UserHandler) handleRegister(ctx *fiber.Ctx) error {
	var req requesttype.RegisterRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	{
		if !model.IsValidUsername(req.Name) {
			return ctx.Status(http.StatusBadRequest).JSON(responsetype.RegisterFail{Message: "Username hanya boleh mengandung huruf, angka, dan underscore. Panjang username minimal 3 karakter dan maksimal 20 karakter"})
		}
	}

	otp := req.OTP

	if err := h.otpService.VerifyOTP(otp.ID, req.Email, model.OTPScopeRegister, otp.Token, otp.Code); err != nil {
		switch {
		case errors.Is(err, service.ErrOTPInvalidCode):
			return errortype.ErrInvalidOTPCode
		case errors.Is(err, service.ErrOTPInvalidScope), errors.Is(err, service.ErrOTPInvalidToken), errors.Is(err, service.ErrOTPAlreadyUsed), errors.Is(err, service.ErrOTPExpired), errors.Is(err, database.ErrNoRows):
			return errortype.ErrInvalidOTP
		default:
			return err
		}
	}

	user, err := h.userService.CreateUser(model.CreateUser{
		Name:           req.Name,
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		Password:       req.Password,
		DeviceUniqueID: req.DeviceUniqueID,
	})
	if err != nil {
		return err
	}

	rawTok := random.Token()
	token, err := h.authService.CreateToken(user, rawTok)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(responsetype.RegisterSuccess{Token: token.Format(rawTok)})
}

func (h *UserHandler) handleRegisterWithGoogle(ctx *fiber.Ctx) error {
	var req requesttype.RegisterWithGoogleRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	info, err := h.authService.ValidateGoogleToken(req.IDToken)
	if err != nil {
		return err
	}
	user, err := h.userService.GetUserByEmail(info.Email)
	if err == nil {
		// login, for better UX, some users might not remember that they have registered before
		rawTok := random.Token()
		token, err := h.authService.CreateToken(user, rawTok)
		if err != nil {
			return err
		}
		return ctx.Status(http.StatusOK).JSON(responsetype.LoginSuccess{Token: token.Format(rawTok)})
	}
	if errors.Is(err, database.ErrNoRows) {
		generatedUsername, err := h.userService.GenerateUsernameFromGoogleAuth(info)
		if err != nil {
			return err
		}
		displayName := helper.StringRemoveNonAlphanumericWithoutSpace(info.DisplayName)
		if len(displayName) >= 25 {
			displayName = displayName[:25]
		}
		if len(displayName) < 1 {
			displayName = "User"
		}
		user2, err := h.userService.CreateUser(model.CreateUser{
			Name:           generatedUsername,
			Email:          info.Email,
			DisplayName:    displayName,
			Password:       random.String(32), // TODO: hacky, but we need to set password to something
			DeviceUniqueID: req.DeviceUniqueID,
		})
		if err != nil {
			return err
		}
		rawTok := random.Token()
		token, err := h.authService.CreateToken(user2, rawTok)
		if err != nil {
			return err
		}
		return ctx.Status(http.StatusOK).JSON(responsetype.RegisterSuccess{Token: token.Format(rawTok)})
	}
	return err
}
