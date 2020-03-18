package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	// подумать насчёт конфига из вайпера, чтобы оттуда тягать адрес фронта
)

// инициализировать поля структуры с конфига
type GoMiddleware struct {
}

func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}

// TODO: убрать хардкод
func (mw *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:5757")
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
			}
		}()
		return next(c)
	}
}
