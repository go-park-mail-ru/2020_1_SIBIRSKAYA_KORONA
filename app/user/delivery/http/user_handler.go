package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	useCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase) {
	handler := &UserHandler{
		useCase: useCase,
	}
	router.POST("/settings", handler.Create)
	router.GET("/profile/:user", handler.Get) // по id или nicName
	router.GET("/settings", handler.GetAll)   // получ все настройки
	router.PUT("/settings", handler.Update)
	router.DELETE("/settings", handler.Delete)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	// в миддлвар
	if _, err := ctx.Cookie("session_id"); err == nil {
		return ctx.NoContent(http.StatusSeeOther)
	}
	//

	reqBody, err := ioutil.ReadAll(ctx.Request().Body)
	usr := models.Create(reqBody)
	if err != nil ||  usr == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()
	sid, err := userHandler.useCase.Create(usr)
	if err != nil {
		return ctx.NoContent(http.StatusConflict)
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sid,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	ctx.SetCookie(cookie)
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Get(ctx echo.Context) error {
	userData := userHandler.useCase.Get(ctx.Param("user"))
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
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	userData := userHandler.useCase.GetAll(cookie.Value)
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
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	newUser := new(models.User)
	newUser.Name = ctx.FormValue("newName")
	newUser.Surname = ctx.FormValue("newSurname")
	newUser.Nickname= ctx.FormValue("newNickname")
	newUser.Email = ctx.FormValue("newEmail")
	newUser.Password = ctx.FormValue("newPassword")
	oldPass := ctx.FormValue("oldPassword")

	/*
	* file, err := ctx.FormFile("avatar")
	* if err != nil {
	*     return err
	* }
	*/

	// TODO класс ошибок
	if userHandler.useCase.Update(cookie.Value, oldPass, newUser) != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Delete(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	if userHandler.useCase.Delete(cookie.Value) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}
