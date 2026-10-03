package handler

import (
	"errors"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/password"
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
	"strings"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Route(r *fiber.App) {
	group := r.Group("/auth")
	group.Post("/login", h.handleLogin)
	group.Post("/loginWithGoogle", h.handleLoginWithGoogle)

	auth := middleware.AuthMiddleware(h.authService)
	group.Get("/details", auth, h.handleDetails)
	group.Post("/logout", auth, h.handleLogout)
	group.Post("/extendToken", auth, h.handleExtendToken)
}

func (h *AuthHandler) handleLoginWithGoogle(ctx *fiber.Ctx) error {
	var req requesttype.LoginWithGoogleRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	payload, err := h.authService.ValidateGoogleToken(req.IDToken)
	if err != nil {
		return err
	}
	user, err := h.userService.GetUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return errortype.ErrEmailNotRegistered
		}
		return err
	}
	rawTok := random.Token()
	tok, err := h.authService.CreateToken(user, rawTok)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(&responsetype.LoginSuccess{
		Token: tok.Format(rawTok),
	})
}

func (h *AuthHandler) handleExtendToken(ctx *fiber.Ctx) error {
	if err := h.authService.ExtendToken(ctxhelper.GetToken(ctx)); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Token extended"))
}

func (h *AuthHandler) getUserByNameOrEmail(emailOrUsername string) (model.User, error) {
	if strings.Contains(emailOrUsername, "@") {
		return h.userService.GetUserByEmail(emailOrUsername)
	}
	return h.userService.GetUserByName(emailOrUsername)
}

func (h *AuthHandler) handleLogin(ctx *fiber.Ctx) error {
	var req requesttype.LoginRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	user, err := h.getUserByNameOrEmail(strings.ToLower(req.Username))
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return err
		}
		return errortype.ErrInvalidUsernameOrPassword
	}

	if err := password.Verify(user.PasswordHash, []byte(req.Password)); err != nil {
		return errortype.ErrInvalidUsernameOrPassword
	}

	rawTok := random.Token()
	tok, err := h.authService.CreateToken(user, rawTok)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(&responsetype.LoginSuccess{
		Token: tok.Format(rawTok),
	})
}

func (h *AuthHandler) handleDetails(ctx *fiber.Ctx) error {
	token := ctxhelper.GetToken(ctx)

	user, err := h.userService.GetUser(token.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(&responsetype.TokenDetails{
		UserID:          token.UserID,
		UserName:        user.Name,
		UserDisplayName: user.DisplayName,
		UserEmail:       user.Email,
		UserRole:        token.Role,
		HasPin:          user.PinHash != nil,
	})
}

func (h *AuthHandler) handleLogout(ctx *fiber.Ctx) error {
	token := ctxhelper.GetToken(ctx)

	if err := h.authService.RevokeToken(token.ID); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK)
}
