package ctxhelper

import (
	"errors"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func ReadFiles(ctx *fiber.Ctx) ([]*model.File, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		if err == fasthttp.ErrNoMultipartForm {
			return []*model.File{}, nil
		}
		return nil, err
	}

	files := make([]*model.File, 0)
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			if fileHeader.Size >= 1e8 {
				return nil, errors.New("file must under 100mb")
			}
			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}
			files = append(files, &model.File{Name: fileHeader.Filename, Bytes: file})
		}
	}
	return files, nil
}
