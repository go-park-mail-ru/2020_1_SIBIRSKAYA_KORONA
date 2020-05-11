package middleware

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/csrf"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/metric"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/sanitize"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/webSocketPool"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Middleware struct {
	origins    map[string]struct{}
	serverMode string

	metr   metric.Metrics
	wsPool webSocketPool.WebSocketPool

	sUseCase    session.UseCase
	uUseCase    user.UseCase
	bUseCase    board.UseCase
	cUseCase    column.UseCase
	tUseCase    task.UseCase
	comUseCase  comment.UseCase
	lUseCase    label.UseCase
	chUseCase   checklist.UseCase
	itUseCase   item.UseCase
	atchUseCase attach.UseCase
	notfUseCase notification.UseCase
}

func CreateMiddleware(sUseCase_ session.UseCase, uUseCase_ user.UseCase, bUseCase_ board.UseCase, cUseCase_ column.UseCase,
	tUseCase_ task.UseCase, comUseCase_ comment.UseCase, chUseCase_ checklist.UseCase, itUseCase_ item.UseCase,
	lUseCase_ label.UseCase, atchUseCase_ attach.UseCase, notfUseCase_ notification.UseCase) *Middleware {
	origins_ := make(map[string]struct{})
	// TODO: вайпер
	for _, key := range viper.GetStringSlice("cors.allowed_origins") {
		origins_[key] = struct{}{}
	}
	return &Middleware{
		origins:    origins_,
		serverMode: viper.GetString("server.mode"),

		sUseCase:    sUseCase_,
		uUseCase:    uUseCase_,
		bUseCase:    bUseCase_,
		cUseCase:    cUseCase_,
		lUseCase:    lUseCase_,
		tUseCase:    tUseCase_,
		comUseCase:  comUseCase_,
		chUseCase:   chUseCase_,
		itUseCase:   itUseCase_,
		atchUseCase: atchUseCase_,
		notfUseCase: notfUseCase_,
	}
}

func (mw *Middleware) SetMetrics(metr_ metric.Metrics) {
	mw.metr = metr_
}

func (mw *Middleware) SetWebSocketPool(wsPool_ webSocketPool.WebSocketPool) {
	mw.wsPool = wsPool_
}

func (mw *Middleware) RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
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

func (mw *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		origin := ctx.Request().Header.Get("Origin")
		if _, exist := mw.origins[origin]; !exist {
			return ctx.NoContent(http.StatusForbidden)
		}
		ctx.Response().Header().Set("Access-Control-Allow-Origin", origin)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		if ctx.Request().Method == "OPTIONS" {
			return ctx.NoContent(http.StatusOK)
		}
		return next(ctx)
	}
}

func (mw *Middleware) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
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

func (mw *Middleware) Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)
		var status int
		if err != nil {
			status = err.(*echo.HTTPError).Code
		} else {
			status = ctx.Response().Status
		}
		mw.metr.ObserveResponseTime(status, ctx.Request().Method, ctx.Path(), time.Since(start).Seconds())
		mw.metr.IncHits(status, ctx.Request().Method, ctx.Path())
		return nil
	}
}

func (mw *Middleware) Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer ctx.Request().Body.Close()
		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("body", sanBody)
		return next(ctx)
	}
}

func (mw *Middleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			return ctx.String(http.StatusUnauthorized, errors.NoCookie)
		}
		sid := cookie.Value
		uid, has := mw.sUseCase.Get(sid)
		if !has {
			// Пришла невалидная кука, стираем её из браузера
			newCookie := http.Cookie{Name: "session_id", Value: sid, Expires: time.Now().AddDate(-1, 0, 0)}
			ctx.SetCookie(&newCookie)
			return ctx.String(http.StatusUnauthorized, errors.NoCookie)
		}
		ctx.Set("uid", uid)
		ctx.Set("sid", sid)
		return next(ctx)
	}
}

func (mw *Middleware) CSRFmiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get(csrf.CSRFheader)
		if token == "" {
			return ctx.String(http.StatusForbidden, errors.DetectedCSRF)
		}
		sid := ctx.Get("sid").(string)
		if !csrf.ValidateToken(token, sid) {
			return ctx.String(http.StatusForbidden, errors.DetectedCSRF)
		}
		return next(ctx)
	}
}

// Вызывается после AuthByCookie
func (mw *Middleware) CheckBoardMemberPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var bid uint
		_, err := fmt.Sscan(ctx.Param("bid"), &bid)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		uid := ctx.Get("uid").(uint)
		if _, err := mw.bUseCase.Get(uid, bid, false); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("bid", bid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckBoardAdminPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var bid uint
		_, err := fmt.Sscan(ctx.Param("bid"), &bid)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		uid := ctx.Get("uid").(uint)
		if _, err := mw.bUseCase.Get(uid, bid, true); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("bid", bid)
		return next(ctx)
	}
}

// вызывается после CheckBoard...Permission
func (mw *Middleware) CheckUserForAssignInBoard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		bid := ctx.Get("bid").(uint)
		var assignUid uint
		_, err := fmt.Sscan(ctx.Param("uid"), &assignUid)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.bUseCase.Get(assignUid, bid, false); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("forUid", assignUid)
		ctx.Set("bid", bid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckColInBoard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		bid := ctx.Get("bid").(uint)
		var cid uint
		if _, err := fmt.Sscan(ctx.Param("cid"), &cid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.cUseCase.Get(bid, cid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("cid", cid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckLabelInBoard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		bid := ctx.Get("bid").(uint)
		var lid uint
		if _, err := fmt.Sscan(ctx.Param("lid"), &lid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.lUseCase.Get(bid, lid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("lid", lid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckTaskInCol(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cid := ctx.Get("cid").(uint)
		var tid uint
		if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.tUseCase.Get(cid, tid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("tid", tid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckCommentInTask(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tid := ctx.Get("tid").(uint)
		var comid uint
		if _, err := fmt.Sscan(ctx.Param("comid"), &comid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.comUseCase.GetByID(tid, comid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("comid", comid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckChecklistInTask(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tid := ctx.Get("tid").(uint)
		var clid uint
		if _, err := fmt.Sscan(ctx.Param("clid"), &clid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.chUseCase.GetByID(tid, clid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("clid", clid)
		return next(ctx)

	}
}

func (mw *Middleware) CheckItemInChecklist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		clid := ctx.Get("clid").(uint)
		var itid uint
		if _, err := fmt.Sscan(ctx.Param("itid"), &itid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.itUseCase.GetByID(clid, itid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("itid", itid)
		return next(ctx)
	}
}

func (mw *Middleware) CheckAttachInTask(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tid := ctx.Get("tid").(uint)
		var fid uint
		if _, err := fmt.Sscan(ctx.Param("fid"), &fid); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		if _, err := mw.atchUseCase.GetByID(tid, fid); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("fid", fid)
		return next(ctx)
	}
}

func (mw *Middleware) DebugMiddle(next echo.HandlerFunc) echo.HandlerFunc {
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

// notification
func (mw *Middleware) SendInviteNotifications(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		status := ctx.Response().Status
		if err != nil || status != http.StatusOK {
			logger.Error("error:", err, " status:", status)
			return err
		}
		ev := &models.Event{
			EventType: ctx.Get("eventType").(string),
			IsRead:    false,
			MakeUid:   ctx.Get("uid").(uint),
		}
		ev.MakeUsr, err = mw.uUseCase.GetByID(ev.Uid)
		if err != nil {
			logger.Error(err)
			return err
		}
		// TODO: вынести в отдельные функции
		var membes models.Users
		if ev.EventType == "InviteToBoard" {
			ev.MetaData.Uid = ctx.Get("forUid").(uint)
			ev.MetaData.Usr, err = mw.uUseCase.GetByID(ev.MetaData.Uid)
			if err != nil {
				logger.Error(err)
				return err
			}
			ev.MetaData.Bid = ctx.Get("bid").(uint)
			tmp, err := mw.bUseCase.Get(ev.MakeUid, ev.MetaData.Bid, false)
			if err != nil {
				logger.Error(err)
				return err
			}
			membes = append(tmp.Members, tmp.Admins...)
		} else if ev.EventType == "AssignOnTask" {
			ev.MetaData.Uid = ctx.Get("forUid").(uint)
			ev.MetaData.Bid = ctx.Get("bid").(uint)
			ev.MetaData.Cid = ctx.Get("cid").(uint)
			ev.MetaData.Tid = ctx.Get("tid").(uint)
			tmp, err := mw.tUseCase.Get(ev.MetaData.Cid, ev.MetaData.Tid)
			if err != nil {
				logger.Error(err)
				return err
			}
			membes = tmp.Members
		}
		for _, elem := range membes {
			ev.Uid = elem.ID
			if err = mw.notfUseCase.Create(ev); err != nil {
				logger.Error(err)
			}
			resp, err := ev.MarshalJSON()
			if err != nil {
				logger.Error(err)
			}
			mw.wsPool.Send(ev.Uid, resp)
			logger.Info("send notifications to user:", ev.Uid)
		}
		return nil
	}
}
