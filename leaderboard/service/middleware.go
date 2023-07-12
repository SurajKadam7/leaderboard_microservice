package leaderservice

import (
	"context"

	"github.com/go-kit/log"
	"github.com/surajkadam/youtube_assignment/model"
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

func (mw loggingMiddleware) DayTopViewd(ctx context.Context, limit int64) (result []model.Response, err error) {
	defer func() {
		mw.logger.Log("method", "DayTopViewd", "limit", limit, "err", err)
	}()
	return mw.next.DayTopViewd(ctx, limit)
}

func (mw loggingMiddleware) LifetimeTopViews(ctx context.Context, limit int64) (result []model.Response, err error) {
	defer func() {
		mw.logger.Log("method", "LifetimeTopViews", "limit", limit, "err", err)
	}()
	return mw.next.LifetimeTopViews(ctx, limit)
}
