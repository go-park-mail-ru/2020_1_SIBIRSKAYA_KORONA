package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type GoMiddleware struct {
	frontendUrl string
	serverMode  string

	sRepo    session.Repository
	bUsecase board.UseCase
}

func InitMiddleware(sRepo_ session.Repository, bUsecase_ board.UseCase) *GoMiddleware {
	return &GoMiddleware{
		frontendUrl: fmt.Sprintf("%s://%s:%s",
			viper.GetString("frontend.protocol"),
			viper.GetString("frontend.ip"),
			viper.GetString("frontend.port")),

		serverMode: viper.GetString("server.mode"),

		sRepo:    sRepo_,
		bUsecase: bUsecase_,
	}
}

func (mw *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Origin", mw.frontendUrl)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return next(ctx)
	}
}

func (mw *GoMiddleware) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() error {
			if err := recover(); err != nil {
				fmt.Println("ProcessPanic up on ", ctx.Request().Method, ctx.Request().URL.Path)
				fmt.Println("Panic statement: ", err)
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return nil
		}()
		return next(ctx)
	}
}

func (mw *GoMiddleware) AuthByCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			return ctx.NoContent(http.StatusForbidden)
		}
		sid := cookie.Value
		userID, exist := mw.sRepo.Get(sid)
		if exist != true {
			return ctx.NoContent(http.StatusNotFound)
		}
		ctx.Set("userID", userID)
		ctx.Set("sessionID", sid)
		return next(ctx)
	}
}

// Вызывается после AuthByCookie
func (mw *GoMiddleware) CheckBoardMemberPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var bid uint
		_, err := fmt.Sscan(ctx.Param("bid"), &bid)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		uid := ctx.Get("userID").(uint)
		if _, err := mw.bUsecase.Get(uid, bid, false); err != nil {
			return ctx.NoContent(http.StatusUnauthorized)
		}
		return next(ctx)
	}
}

func (mw *GoMiddleware) CheckBoardAdminPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var bid uint
		_, err := fmt.Sscan(ctx.Param("bid"), &bid)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		uid := ctx.Get("userID").(uint)

		if _, err := mw.bUsecase.Get(uid, bid, true); err != nil {
			return ctx.NoContent(http.StatusUnauthorized)
		}
		return next(ctx)
	}
}

func (mw *GoMiddleware) DebugMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if mw.serverMode == "debug" {
			dump, err := httputil.DumpRequest(ctx.Request(), true)
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			logger.Debugf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump)
		}

		return next(ctx)
	}
}
