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
	router.GET("/settings", handler.Join)
}

// TODO: мидлвары на валидацию, запрос куки, панику, ошибку
func (sessionHandler *SessionHandler) Join(ctx echo.Context) error {
	u := new(models.User)
	if err := ctx.Bind(u); err != nil {
		return err
	}
	sid, err := sessionHandler.useCase.Create(u)
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
func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}

func (sessionHandler *SessionHandler) LogOut(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}

func (sessionHandler *SessionHandler) Delete(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}