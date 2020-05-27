package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/template"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/labstack/echo/v4"
)

type TemplateHandler struct {
	usecase template.Usecase
}

func CreateHandler(router *echo.Echo, usecase_ template.Usecase, mw *middleware.Middleware) {
	handler := &TemplateHandler{
		usecase: usecase_,
	}

	router.POST("/boards/templates/:name", handler.Create, mw.Sanitize, mw.CheckAuth)
}

func (templateHandler *TemplateHandler) Create(ctx echo.Context) error {
	var tmpl models.Template
	body := ctx.Get("body").([]byte)
	err := tmpl.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	uid := ctx.Get("uid").(uint)
	err = templateHandler.usecase.Create(uid, &tmpl)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	// Что возвращать будем ?
	resp, err := tmpl.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
