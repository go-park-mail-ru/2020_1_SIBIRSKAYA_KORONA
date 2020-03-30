package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type GoMiddleware struct {
	frontendUrl string
}

func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{
		frontendUrl: fmt.Sprintf("%s://%s:%s",
			viper.GetString("frontend.protocol"),
			viper.GetString("frontend.ip"),
			viper.GetString("frontend.port")),
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

func (mw *GoMiddleware) CheckCookieExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			return ctx.NoContent(http.StatusForbidden)
		}
		ctx.Set("sid", cookie.Value)
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
