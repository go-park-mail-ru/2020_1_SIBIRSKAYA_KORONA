package http

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/csrf"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	UseCase session.UseCase
}

func CreateHandler(router *echo.Echo, useCase session.UseCase, mw *middleware.Middleware) {
	handler := &SessionHandler{
		UseCase: useCase,
	}
	router.POST("/api/session", handler.LogIn, mw.Sanitize)
	router.GET("/api/token", handler.Token, mw.CheckAuth)
	router.DELETE("/api/session", handler.LogOut, mw.CheckAuth)
}

func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {
	var usr models.User
	body := ctx.Get("body").([]byte)
	err := usr.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	sessionExpires := time.Now().AddDate(1, 0, 0)
	if sid, err := sessionHandler.UseCase.Create(&usr, int32(sessionExpires.Unix())); err != nil {
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	} else {
		cookie := &http.Cookie{
			Name:    "session_id",
			Value:   sid,
			Path:    "/",
			Expires: sessionExpires,
			// SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		}
		ctx.SetCookie(cookie)
		return ctx.NoContent(http.StatusOK)
	}
}

func (sessionHandler *SessionHandler) Token(ctx echo.Context) error {
	sid := ctx.Get("sid").(string)
	token := csrf.MakeToken(sid)
	ctx.Response().Header().Set(csrf.CSRFheader, token)
	ctx.Response().Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
	return ctx.NoContent(http.StatusOK)
}

func (sessionHandler *SessionHandler) LogOut(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)
	err := sessionHandler.UseCase.Delete(cookie)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	newCookie := http.Cookie{Name: "session_id", Value: cookie, Expires: time.Now().AddDate(-1, 0, 0)}
	ctx.SetCookie(&newCookie)
	return ctx.NoContent(http.StatusOK)
}
