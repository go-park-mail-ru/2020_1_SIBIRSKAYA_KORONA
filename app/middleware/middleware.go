package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/csrf"
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
	tUseCase task.UseCase
}

func CreateMiddleware(sUseCase_ session.UseCase, bUseCase_ board.UseCase, cUseCase_ column.UseCase, tUseCase_ task.UseCase) *GoMiddleware {
	return &GoMiddleware{
		frontendUrl: fmt.Sprintf("%s://%s:%s",
			viper.GetString("frontend.protocol"),
			viper.GetString("frontend.ip"),
			viper.GetString("frontend.port")),

		serverMode: viper.GetString("server.mode"),

		sUseCase: sUseCase_,
		bUseCase: bUseCase_,
		cUseCase: cUseCase_,
		tUseCase: tUseCase_,
	}
}

func (mw *GoMiddleware) RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		res := next(ctx)
		logger.Infof("%s %s %d %s",
			ctx.Request().Method,
			ctx.Request().RequestURI,
			ctx.Response().Status,
			time.Since(start))

		return res
	}
}

func (mw *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Request().Header.Get("Origin") != mr.frontendUrl {
			return ctx.NoContent(http.StatusForbidden)
		}

		ctx.Response().Header().Set("Access-Control-Allow-Origin", mw.frontendUrl)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		if ctx.Request().Method == "OPTIONS" {
			return ctx.NoContent(http.StatusOK)
		}
		return next(ctx)
	}
}

func (mw *GoMiddleware) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("ProcessPanic up on ", ctx.Request().Method, ctx.Request().URL.Path)
				fmt.Println("Panic statement: ", err)
				ctx.NoContent(http.StatusInternalServerError)
			}
		}()
		return next(ctx)
	}
}

func (mw *GoMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, message.ResponseError{Message: errors.ErrNoCookie.Error()})
		}
		sid := cookie.Value
		uid, has := mw.sUseCase.Get(sid)
		if !has {
			// Пришла невалидная кука, стираем её из браузера
			newCookie := http.Cookie{Name: "session_id", Value: sid, Expires: time.Now().AddDate(-1, 0, 0)}
			ctx.SetCookie(&newCookie)
			return ctx.JSON(http.StatusUnauthorized, message.ResponseError{Message: errors.ErrNoCookie.Error()})
		}
		ctx.Set("uid", uid)
		ctx.Set("sid", sid)
		return next(ctx)
	}
}

func (mw *GoMiddleware) CSRFmiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get(csrf.CSRFheader)
		if token == "" {
			return ctx.JSON(http.StatusForbidden, message.ResponseError{Message: errors.ErrDetectedCSRF.Error()})
		}

		sid := ctx.Get("sid").(string)

		if !csrf.ValidateToken(token, sid) {
			return ctx.JSON(http.StatusForbidden, message.ResponseError{Message: errors.ErrDetectedCSRF.Error()})
		}

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
		uid := ctx.Get("uid").(uint)
		if _, err := mw.bUseCase.Get(uid, bid, false); err != nil {
			logger.Error(err)
			return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
		}
		ctx.Set("bid", bid)
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
		uid := ctx.Get("uid").(uint)
		if _, err := mw.bUseCase.Get(uid, bid, true); err != nil {
			logger.Error(err)
			return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
		}
		ctx.Set("bid", bid)
		return next(ctx)
	}
}

func (mw *GoMiddleware) CheckColInBoard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		bid := ctx.Get("bid").(uint)
		var cid uint
		if _, err := fmt.Sscan(ctx.Param("cid"), &cid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.cUseCase.Get(bid, cid); err != nil {
			logger.Error(err)
			return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
		}
		ctx.Set("cid", cid)
		return next(ctx)
	}
}

func (mw *GoMiddleware) CheckTaskInCol(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cid := ctx.Get("cid").(uint)
		var tid uint
		if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.tUseCase.Get(cid, tid); err != nil {
			logger.Error(err)
			return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
		}
		ctx.Set("tid", tid)
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
