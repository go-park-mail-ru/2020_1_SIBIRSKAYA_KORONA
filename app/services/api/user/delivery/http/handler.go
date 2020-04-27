package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type UserHandler struct {
	useCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase, mw *middleware.Middleware) {
	handler := &UserHandler{
		useCase: useCase,
	}
	router.POST("/settings", handler.Create, mw.Sanitize)
	router.GET("/profile/:id_or_nickname", handler.Get)
	router.GET("/settings", handler.GetAll, mw.CheckAuth) // получ все настройки
	router.PUT("/settings", handler.Update, mw.CheckAuth, mw.CSRFmiddle)
	router.DELETE("/settings", handler.Delete, mw.CheckAuth)
	router.GET("/search/profile", handler.GetUsersByNicknamePart, mw.CheckAuth)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	var usr models.User
	body := ctx.Get("body").([]byte)
	err := usr.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	// TODO: вайпер
	usr.Avatar = fmt.Sprintf("%s://%s:%s%s",
		viper.GetString("frontend.protocol"),
		viper.GetString("frontend.ip"),
		viper.GetString("frontend.port"),
		viper.GetString("frontend.default_avatar"))
	sessionExpires := time.Now().AddDate(1, 0, 0)
	sid, err := userHandler.useCase.Create(&usr, int32(sessionExpires.Unix()))
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
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
	usrKey := ctx.Param("id_or_nickname")
	usr := new(models.User)
	var err error
	if uid, er := strconv.Atoi(usrKey); er == nil {
		usr, err = userHandler.useCase.GetByID(uint(uid))
	} else {
		usr, err = userHandler.useCase.GetByNickname(usrKey)
	}
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := usr.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (userHandler *UserHandler) GetAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	usr, err := userHandler.useCase.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := usr.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (userHandler *UserHandler) Update(ctx echo.Context) error {
	var newUser models.User
	newUser.ID = ctx.Get("uid").(uint)
	newUser.Name = ctx.FormValue("newName")
	newUser.Surname = ctx.FormValue("newSurname")
	newUser.Nickname = ctx.FormValue("newNickname")
	newUser.Email = ctx.FormValue("newEmail")
	newUser.Password = []byte(ctx.FormValue("newPassword"))
	oldPass := []byte(ctx.FormValue("oldPassword"))
	avatarFileDescriptor, err := ctx.FormFile("avatar")
	if err != nil {
		logger.Error(err)
	}
	if err := userHandler.useCase.Update(oldPass, newUser, avatarFileDescriptor); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Delete(ctx echo.Context) error {
	sid := ctx.Get("sid").(string)
	uid := ctx.Get("uid").(uint)
	if userHandler.useCase.Delete(uid, sid) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	newCookie := http.Cookie{Name: "session_id", Value: sid, Expires: time.Now().AddDate(-1, 0, 0)}
	ctx.SetCookie(&newCookie)
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) GetUsersByNicknamePart(ctx echo.Context) error {
	nicknamePart := ctx.QueryParam("nickname")
	if nicknamePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}
	var limit uint
	_, err := fmt.Sscan(ctx.QueryParam("limit"), &limit)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	_, err = userHandler.useCase.GetUsersByNicknamePart(nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// body, err := message.GetBody(message.Pair{Name: "user", Data: usersData})
	// if err != nil {
	// 	return ctx.NoContent(http.StatusInternalServerError)
	// }
	return ctx.String(http.StatusOK, "мой репозиторий не работает")
}
