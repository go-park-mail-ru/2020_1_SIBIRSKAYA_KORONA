package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type UserHandler struct {
	useCase user.UseCase
}

func CreateHandlerTest(useCase user.UseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func CreateHandler(router *echo.Echo, useCase user.UseCase, mw *middleware.GoMiddleware) {
	handler := &UserHandler{
		useCase: useCase,
	}
	router.POST("/settings", handler.Create)
	router.GET("/profile/:id_or_nickname", handler.Get)
	router.GET("/settings", handler.GetAll, mw.CheckAuth) // получ все настройки
	router.GET("/boards", handler.GetBoards, mw.CheckAuth)
	router.PUT("/settings", handler.Update, mw.CheckAuth, mw.CSRFmiddle)
	router.DELETE("/settings", handler.Delete, mw.CheckAuth)
	//GET       /search/profile?nickname={part_of_nickname}
	router.GET("/search/profile", handler.GetUsersByNicknamePart, mw.CheckAuth)
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

	usersData, err := userHandler.useCase.GetUsersByNicknamePart(nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}

	body, err := message.GetBody(message.Pair{Name: "user", Data: usersData})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	usr := models.CreateUser(ctx)
	if usr == nil {
		log.Println("bad bad bad")
		return ctx.NoContent(http.StatusBadRequest)
	}
	usr.Avatar = fmt.Sprintf("%s://%s:%s%s",
		viper.GetString("frontend.protocol"),
		viper.GetString("frontend.ip"),
		viper.GetString("frontend.port"),
		viper.GetString("frontend.default_avatar"))
	sessionExpires := time.Now().AddDate(1, 0, 0)
	sid, err := userHandler.useCase.Create(usr, sessionExpires.Unix())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
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
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "user", Data: *usr})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) GetAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	userData, err := userHandler.useCase.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "user", Data: *userData})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (userHandler *UserHandler) GetBoards(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	bAdmin, bMember, err := userHandler.useCase.GetBoardsByID(uid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "admin", Data: bAdmin}, message.Pair{Name: "member", Data: bMember})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
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
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
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
