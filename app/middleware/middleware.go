package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// инициализировать поля структуры с конфига
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

// TODO: убрать хардкод
func (mw *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Origin", mw.frontendUrl)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		//c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return next(ctx)
	}
}

// TODO: пришить логгер
func (mw *GoMiddleware) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("panicMiddleware", c.Request().URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Panic is catched: ", err)
				// TODO: решить какой ответ отдавать клиенту
			}
		}()
		return next(c)
	}
}
