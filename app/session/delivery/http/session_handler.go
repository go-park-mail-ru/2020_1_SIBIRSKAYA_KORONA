package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
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

func (sessionHandler *SessionHandler) LogIn(ctx echo.Context) error {
	// в миддлвар
	if _, err := ctx.Cookie("session_id"); err == nil {
		return ctx.NoContent(http.StatusSeeOther)
	}
	//

	reqBody, err := ioutil.ReadAll(ctx.Request().Body)
	usr := models.CreateUser(reqBody)
	if err != nil ||  usr == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()
	sessionExpires := time.Now().AddDate(1,0,0)
	sid, err := sessionHandler.useCase.Create(usr, sessionExpires)
	if err != nil {
		return ctx.NoContent(http.StatusConflict)
	}
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

func (sessionHandler *SessionHandler) IsAuth(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	if !sessionHandler.useCase.Has(cookie.Value) {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (sessionHandler *SessionHandler) LogOut(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	if sessionHandler.useCase.Delete(cookie.Value) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookie.Expires = time.Now().AddDate(-1, 0, 0)
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}