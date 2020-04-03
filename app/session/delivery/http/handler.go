package http

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	useCase session.UseCase
}

type ResponseError struct {
	Message string `json:"message"`
}

func CreateHandler(router *echo.Echo, useCase session.UseCase, mw *middleware.GoMiddleware) {
	handler := &SessionHandler{
		useCase: useCase,
	}
	router.OPTIONS("/session", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

	// TODO: решить как вешать мидлу на handler.LogIn
	router.POST("/session", handler.LogIn)
	router.GET("/session", handler.IsAuth, mw.AuthByCookie)
	router.DELETE("/session", handler.LogOut, mw.AuthByCookie)
}

func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {
	// в миддлвар
	if _, err := ctx.Cookie("session_id"); err == nil {
		return ctx.NoContent(http.StatusSeeOther)
	}
	//

	usr := models.CreateUser(ctx)
	if usr == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()
	sessionExpires := time.Now().AddDate(1, 0, 0)
	if sid, useErr := sessionHandler.useCase.Create(usr, sessionExpires); useErr != nil {
		return ctx.JSON(errors.ResolveErrorToCode(useErr), ResponseError{Message: useErr.Error()})
	} else {
		cookie := &http.Cookie{
			Name:    "session_id",
			Value:   sid,
			Path:    "/",
			Expires: sessionExpires,
			// SameSite: http.SameSiteStrictMode,
		}
		ctx.SetCookie(cookie)
		return ctx.NoContent(http.StatusOK)
	}
}

func (sessionHandler *SessionHandler) IsAuth(ctx echo.Context) error {
	cookie := ctx.Get("sessionID").(string)

	if !sessionHandler.useCase.Has(cookie) {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (sessionHandler *SessionHandler) LogOut(ctx echo.Context) error {
	cookie := ctx.Get("sessionID").(string)

	if sessionHandler.useCase.Delete(cookie) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	newCookie := http.Cookie{Name: "session_id", Value: cookie, Expires: time.Now().AddDate(-1, 0, 0)}
	ctx.SetCookie(&newCookie)

	return ctx.NoContent(http.StatusOK)
}
