package cache

import (
	"context"

	"github.com/surajkadam/youtube_assignment/model"
)

//go:generate mockgen -source=cahe.go -destination mock/mock.go
type Repository interface {
	Viewed(ctx context.Context, video string, incr int64) (res float64, err error)
	DayTopViewed(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error)
	LifetimeTopViewed(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error)
	DayViewCount(ctx context.Context, video string) (viewes float64, err error)
	LifetimeViewCount(ctx context.Context, video string) (viewes float64, err error)
	AddVideos(ctx context.Context, videos []model.Video) (err error)
}
// mockgen . Repository