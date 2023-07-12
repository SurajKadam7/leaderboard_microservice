package leaderservice

import (
	"context"

	"github.com/go-kit/log"

	"github.com/surajkadam/youtube_assignment/config"
	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	"github.com/surajkadam/youtube_assignment/model"
)

type Err struct{}

func (e *Err) Error() string {
	return "bad request"
}

func (e *Err) Status() int {
	return 401
}

type Service interface {
	DayTopViewd(ctx context.Context, limit int64) (result []model.Response, err error)
	LifetimeTopViews(ctx context.Context, limit int64) (result []model.Response, err error)
}

type service struct {
	config *config.Config
}

func New(c *config.Config, logger log.Logger) Service {
	var srv Service
	srv = newService(c)
	srv = LoggingMiddleware(logger)(srv)
	return srv
}

func newService(c *config.Config) *service {
	return &service{
		config: c,
	}
}

func (s *service) top(ctx context.Context, limit int64, key string) (result []model.Response, err error) {
	if limit < 0 {
		return result, youtubeerror.ErrInvalidLimitValue
	}

	if limit > 1000 {
		limit = 1000
	}

	res, err := s.config.Cache.TopViewed(ctx, key, limit)
	if err != nil {
		return nil, err
	}

	return res, nil

}

// need to handle key point..
func (s *service) DayTopViewd(ctx context.Context, limit int64) (result []model.Response, err error) {
	// time.Sleep(time.Second * 10)
	return s.top(ctx, limit, s.config.DayKey())
}

func (s *service) LifetimeTopViews(ctx context.Context, limit int64) (result []model.Response, err error) {
	return s.top(ctx, limit, s.config.LifetimeKey())
}
