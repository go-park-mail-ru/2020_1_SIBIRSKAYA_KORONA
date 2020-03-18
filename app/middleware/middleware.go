package middleware

import (
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

// в echo есть встроенная, но для тренировки пока свои будем писать
// func CORSmiddleware = middleware.CORSWithConfig(middleware.CORSConfig{
// 	AllowOrigins: []string{"http://localhost:5757"},
// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
// })))
