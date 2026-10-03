package ctxhelper

import (
	"errors"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

func doValidate(obj any) error {
	if err := validate.Struct(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]string, 0)
			for _, e := range ve {
				out = append(out, e.Translate(validatorTranslators["id"]))
			}
			return errortype.NewErrorVisible(100000, http.StatusBadRequest, strings.Join(out, "\n"))
		}
		return err
	}
	return nil
}

func BindQuery(ctx *fiber.Ctx, out any) error {
	if err := ctx.QueryParser(out); err != nil {
		return err
	}
	return doValidate(out)
}

func BindBody(ctx *fiber.Ctx, out any) error {
	if err := ctx.BodyParser(out); err != nil {
		return err
	}
	return doValidate(out)
}

func BindHeader(ctx *fiber.Ctx, out any) error {
	if err := ctx.ReqHeaderParser(out); err != nil {
		return err
	}
	return doValidate(out)
}

func BindParam(ctx *fiber.Ctx, out any) error {
	if err := ctx.ParamsParser(out); err != nil {
		return err
	}
	return doValidate(out)
}

func ParseID(ctx *fiber.Ctx) (int64, error) {
	var req requesttype.IDRequest
	if err := BindQuery(ctx, &req); err != nil {
		return 0, err
	}
	return req.ID, nil
}

func ParseParamID(ctx *fiber.Ctx) (int64, error) {
	var req requesttype.IDRequest
	if err := BindParam(ctx, &req); err != nil {
		return 0, err
	}
	return req.ID, nil
}
