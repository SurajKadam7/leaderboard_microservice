package cache

import (
	"context"

	"github.com/surajkadam/youtube_assignment/model"
)

type Cache interface {
	Start() (err error)
	Viewed(ctx context.Context, key string, video string, incr float64) (res float64, err error)
	TopViewed(ctx context.Context, key string, limit int64) (result []model.Response, err error)
	ViewCount(ctx context.Context, key string, video string) (viewes float64, err error)
}
