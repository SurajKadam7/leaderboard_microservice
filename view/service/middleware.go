package viewservice

import (
	"context"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a service Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Viewing(ctx context.Context, video string) (response Response, err error) {
	defer func() {
		mw.logger.Log("method", "Viewing", "video", video, "err", err)
	}()
	return mw.next.Viewing(ctx, video)
}

func (mw loggingMiddleware) DayViews(ctx context.Context, video string) (result Response, err error) {
	defer func() {
		mw.logger.Log("method", "DayViews", "video", video, "err", err)
	}()
	return mw.next.DayViews(ctx, video)
}

func (mw loggingMiddleware) LifetimeViews(ctx context.Context, video string) (result Response, err error) {
	defer func() {
		mw.logger.Log("method", "LifetimeViews", "video", video, "err", err)
	}()
	return mw.next.LifetimeViews(ctx, video)

}
