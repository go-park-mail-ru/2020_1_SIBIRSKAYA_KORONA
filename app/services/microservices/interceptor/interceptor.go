package interceptor

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/metric"

	"google.golang.org/grpc"
)

type Interceptor struct {
	metr metric.Metrics
}

func CreateInterceptor(metr_ metric.Metrics) *Interceptor {
	return &Interceptor{metr: metr_}
}

func (mw *Interceptor) Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	var status = http.StatusOK
	if err != nil {
		status = errors.ResolveErrorToCode(err)
	}
	mw.metr.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	mw.metr.IncHits(status, info.FullMethod, info.FullMethod)
	return resp, err
}
