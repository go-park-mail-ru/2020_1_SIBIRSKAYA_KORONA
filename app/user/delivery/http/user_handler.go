package http

import (
	message "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	useCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase) {
	handler := &UserHandler{
		useCase: useCase,
	}
	router.GET("/profile/:id", handler.Get)
	router.GET("/settings", handler.GetAll) // получ все настройки
	router.POST("/settings", handler.Update)
}

// TODO: мидлвары на валидацию, запрос куки, панику, ошибку
func (userHandler *UserHandler) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	userData := userHandler.useCase.Get(uint(id))
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
	return ctx.String(http.StatusOK, "доделай меня :(")
}

func (userHandler *UserHandler) Update(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}