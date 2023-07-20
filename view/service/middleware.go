package viewservice

import (
	"context"
	"fmt"

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

func (mw loggingMiddleware) View(ctx context.Context, video string) (response model.ViedeoDetails, err error) {
	defer func() {
		mw.logger.Log("method", "Viewing", "video", video, "result", fmt.Sprintf("res : %+v", response), "err", err)
	}()
	return mw.next.View(ctx, video)
}

func (mw loggingMiddleware) DayViews(ctx context.Context, video string) (result model.ViedeoDetails, err error) {
	defer func() {
		mw.logger.Log("method", "DayViews", "video", video, "result", fmt.Sprintf("res : %+v", result), "err", err)
	}()
	return mw.next.DayViews(ctx, video)
}

func (mw loggingMiddleware) LifetimeViews(ctx context.Context, video string) (result model.ViedeoDetails, err error) {
	defer func() {
		mw.logger.Log("method", "LifetimeViews", "video", video, "result", fmt.Sprintf("res : %+v", result), "err", err)
	}()
	return mw.next.LifetimeViews(ctx, video)

}

func (mw loggingMiddleware) AddVideos(ctx context.Context, videos []model.Video) (result model.AddVideoStatus, err error) {
	defer func() {
		mw.logger.Log("method", "AddVideos", "result", fmt.Sprintf("res : %+v", result), "err", err)
	}()
	return mw.next.AddVideos(ctx, videos)

}
