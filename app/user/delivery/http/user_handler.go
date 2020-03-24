package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/custom_errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	useCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase, mw *middleware.GoMiddleware) {
	handler := &UserHandler{
		useCase: useCase,
	}
	router.OPTIONS("/settings", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

	// TODO: решить как вешать мидлу на handler.Create
	router.POST("/settings", handler.Create)
	router.GET("/profile/:user", handler.Get)                    // по id или nicName
	router.GET("/settings", handler.GetAll, mw.CheckCookieExist) // получ все настройки
	router.PUT("/settings", handler.Update, mw.CheckCookieExist)
	router.DELETE("/settings", handler.Delete, mw.CheckCookieExist)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	// в миддлвар
	if _, err := ctx.Cookie("session_id"); err == nil {
		return ctx.NoContent(http.StatusSeeOther)
	}
	//

	reqBody, err := ioutil.ReadAll(ctx.Request().Body)
	usr := models.CreateUser(reqBody)
	if err != nil || usr == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()
	sessionExpires := time.Now().AddDate(1, 0, 0)
	sid, err := userHandler.useCase.Create(usr, sessionExpires)
	if err != nil {
		return ctx.NoContent(http.StatusConflict)
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sid,
		Path:    "/",
		Expires: sessionExpires,
	}
	ctx.SetCookie(cookie)
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Get(ctx echo.Context) error {
	userData := userHandler.useCase.GetByUserKey(ctx.Param("user"))
	if userData == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	body, err := message.GetBody(message.Pair{Name: "user", Data: *userData})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) GetAll(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)

	userData := userHandler.useCase.GetByCookie(cookie)
	if userData == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	body, err := message.GetBody(message.Pair{Name: "user", Data: *userData})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) Update(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)

	newUser := new(models.User)
	newUser.Name = ctx.FormValue("newName")
	newUser.Surname = ctx.FormValue("newSurname")
	newUser.Nickname = ctx.FormValue("newNickname")
	newUser.Email = ctx.FormValue("newEmail")
	newUser.Password = ctx.FormValue("newPassword")
	oldPass := ctx.FormValue("oldPassword")

	/*
	* file, err := ctx.FormFile("avatar")
	* if err != nil {
	*     return err
	* }
	 */

	if err := userHandler.useCase.Update(cookie, oldPass, newUser); err != nil {
		err := err.(*custom_errors.CustomUsecaseError)
		return ctx.String(err.Code, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Delete(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)

	if userHandler.useCase.Delete(cookie) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	newCookie := http.Cookie{Name: "session_id", Value: cookie, Expires: time.Now().AddDate(-1, 0, 0)}
	ctx.SetCookie(&newCookie)

	return ctx.NoContent(http.StatusOK)
}
