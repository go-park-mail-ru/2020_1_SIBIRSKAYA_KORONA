package middleware

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"
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

	sUseCase session.UseCase
	bUseCase board.UseCase
	cUseCase column.UseCase
}

func CreateMiddleware(sUseCase_ session.UseCase, bUseCase_ board.UseCase, cUseCase_ column.UseCase) *GoMiddleware {
	return &GoMiddleware{
		frontendUrl: fmt.Sprintf("%s://%s:%s",
			viper.GetString("frontend.protocol"),
			viper.GetString("frontend.ip"),
			viper.GetString("frontend.port")),

		serverMode: viper.GetString("server.mode"),

		sUseCase: sUseCase_,
		bUseCase: bUseCase_,
		cUseCase: cUseCase_,
	}
}

func (mw *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Origin", mw.frontendUrl)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if ctx.Request().Method == "OPTIONS" {
			return ctx.NoContent(http.StatusOK)
		}
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

func (mw *GoMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			return ctx.NoContent(http.StatusForbidden)
		}
		sid := cookie.Value
		userID, exist := mw.sUseCase.Get(sid)
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
		if _, err := mw.bUseCase.Get(uid, bid, false); err != nil {
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
		if _, err := mw.bUseCase.Get(uid, bid, true); err != nil {
			return ctx.NoContent(http.StatusUnauthorized)
		}
		return next(ctx)
	}
}

func (mw *GoMiddleware) CheckColInBoard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var bid, cid uint
		if _, err := fmt.Sscan(ctx.Param("bid"), &bid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := fmt.Sscan(ctx.Param("cid"), &cid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, useErr := mw.cUseCase.Get(bid, cid); useErr != nil {
			return ctx.JSON(errors.ResolveErrorToCode(useErr), message.ResponseError{Message: useErr.Error()})
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
