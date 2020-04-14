package http

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/csrf"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	useCase session.UseCase
}

func CreateHandlerTest(sessionCase session.UseCase) *SessionHandler {
	return &SessionHandler{
		useCase: sessionCase,
	}
}

func CreateHandler(router *echo.Echo, useCase session.UseCase, mw *middleware.GoMiddleware) {
	handler := &SessionHandler{
		useCase: useCase,
	}
	router.POST("/session", handler.LogIn)
	router.GET("/token", handler.Token, mw.CheckAuth)
	router.DELETE("/session", handler.LogOut, mw.CheckAuth)
}

func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {
	usr := models.CreateUser(ctx)
	if usr == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()
	sessionExpires := time.Now().AddDate(1, 0, 0)
	if sid, err := sessionHandler.useCase.Create(usr, sessionExpires); err != nil {
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
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
	if sessionHandler.useCase.Delete(cookie) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	newCookie := http.Cookie{Name: "session_id", Value: cookie, Expires: time.Now().AddDate(-1, 0, 0)}
	ctx.SetCookie(&newCookie)
	return ctx.NoContent(http.StatusOK)
}
