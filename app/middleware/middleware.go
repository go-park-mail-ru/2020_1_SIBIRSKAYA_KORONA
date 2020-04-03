package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type GoMiddleware struct {
	frontendUrl string
	serverMode  string
}

func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{
		frontendUrl: fmt.Sprintf("%s://%s:%s",
			viper.GetString("frontend.protocol"),
			viper.GetString("frontend.ip"),
			viper.GetString("frontend.port")),

		serverMode: viper.GetString("server.mode"),
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

// TODO: уровень Info
// func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, r)

// 		fmt.Printf("FMT [%s] %s, %s %s\n",
// 			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

// 		log.Printf("LOG [%s] %s, %s %s\n",
// 			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

// 		ac.StdLogger.Printf("[%s] %s, %s %s\n",
// 			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

// 		ac.ZapLogger.Info(r.URL.Path,
// 			zap.String("method", r.Method),
// 			zap.String("remote_addr", r.RemoteAddr),
// 			zap.String("url", r.URL.Path),
// 			zap.Duration("work_time", time.Since(start)),
// 		)

// 		ac.LogrusLogger.WithFields(logrus.Fields{
// 			"method":      r.Method,
// 			"remote_addr": r.RemoteAddr,
// 			"work_time":   time.Since(start),
// 		}).Info(r.URL.Path)
// 	})
// }

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
