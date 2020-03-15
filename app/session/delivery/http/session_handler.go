package http

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	message "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	useCase session.UseCase
}

func CreateHandler(router *echo.Echo, useCase session.UseCase) {
	handler := &SessionHandler{
		useCase: useCase,
	}
	router.POST("/session", handler.LogIn)
	router.GET("/session", handler.IsAuth)
	router.DELETE("/session", handler.LogOut)
}

// TODO: мидлвары на валидацию, запрос куки, панику, ошибку
func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {
	usr := new(models.User)
	if err := ctx.Bind(usr); err != nil {
		return err
	}
	sid, err := sessionHandler.useCase.Create(usr)
	if err != nil {
		body, err := message.GetBody(http.StatusConflict)
		if err != nil {
			return err
		}
		return ctx.String(http.StatusOK, body)
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sid,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
		// SameSite: http.SameSiteStrictMode,
	}
	ctx.SetCookie(cookie)
	body, err := message.GetBody(http.StatusOK)
	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, body)
}

func (sessionHandler *SessionHandler) IsAuth(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}

func (sessionHandler *SessionHandler) LogOut(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return err
	}
	if sessionHandler.useCase.Delete(cookie.Value) != nil {
		return ctx.String(http.StatusOK, "не ок")
	}
	return ctx.String(http.StatusOK, "ок")
}