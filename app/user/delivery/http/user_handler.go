package http

import (
	"encoding/json"
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
	usr := new(models.User)

	// if err := ctx.Bind(usr); err != nil {
	// 	log.Println("111111111111111111")
	// 	return err
	// }

	// ------- костыль от ctx.Bind
	reqBody, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	defer ctx.Request().Body.Close()
	err = json.Unmarshal(reqBody, &usr)
	// -------

	sid, err := userHandler.useCase.Create(usr)
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

func (userHandler *UserHandler) Get(ctx echo.Context) error {
	userData := userHandler.useCase.Get(ctx.Param("user"))
	if userData == nil {
		body, err := message.GetBody(http.StatusNotFound)
		if err != nil {
			return err
		}
		return ctx.String(http.StatusOK, body)
	}
	body, err := message.GetBody(http.StatusOK, message.Pair{Name: "user", Data: *userData})
	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) GetAll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return err
	}
	userData := userHandler.useCase.GetAll(cookie.Value)
	if userData == nil {
		body, err := message.GetBody(http.StatusNotFound)
		if err != nil {
			return err
		}
		return ctx.String(http.StatusOK, body)
	}
	body, err := message.GetBody(http.StatusOK, message.Pair{Name: "user", Data: *userData})
	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) Update(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}

func (userHandler *UserHandler) Delete(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}
