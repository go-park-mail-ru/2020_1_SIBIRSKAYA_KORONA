package http

import (

	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type SessionHandler struct {
	useCase session.UseCase
}

/*func CreateSessionHandler(router *echo.Router, useCase session.UseCase) {
	handler := &SessionHandler{
		useCase: useCase,
	}
	router.GET("/api/user/:nickname/profile", handler.Get)
}*/

// TODO: мидлвары на валидацию, запрос куки, панику
func (sessionHandler *SessionHandler) Join(ctx echo.Context) error {

}

func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {

}

func (sessionHandler *SessionHandler) LogOut(ctx echo.Context) error {

}

func (sessionHandler *SessionHandler) Delete(ctx echo.Context) error {

}