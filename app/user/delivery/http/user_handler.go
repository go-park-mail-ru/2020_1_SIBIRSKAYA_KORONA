package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserHandler struct {
	useCase user.UseCase
}

/*func CreateUserHandler(router *echo.Router, useCase user.UseCase) {
	handler := &UserHandler{
		useCase: useCase,
	}
	router.GET("/api/user/:nickname/profile", handler.Get)
}*/

// TODO: мидлвары на валидацию, запрос куки, панику
func (userHandler *UserHandler) Get(ctx echo.Context) error {
	id := ctx.Param("id")
	userData := userHandler.useCase.Get(id)
	if userData != nil {
		return ctx.JSON(http.StatusOK, "хуево")
	}
	return ctx.JSON(http.StatusOK, userData)
}

func (userHandler *UserHandler) GetAll(ctx echo.Context) error {
	id := ctx.Param("id")
	userData := userHandler.useCase.GetAll(id)
	if userData != nil {
		return ctx.JSON(http.StatusOK, "хуево")
	}
	return ctx.JSON(http.StatusOK, userData)
}

// TODO
func (userHandler *UserHandler) Update(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "доделай меня :(")
}