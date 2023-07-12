package viewservice

import (
	"context"

	"github.com/go-kit/log"
	"github.com/surajkadam/youtube_assignment/config"
)

// my bussiness logic will stay here ...
type Service interface {
	Viewing(ctx context.Context, video string) (result Response, err error)
	DayViews(ctx context.Context, video string) (result Response, err error)
	LifetimeViews(ctx context.Context, video string) (result Response, err error)
}

type Response struct {
	VideoName string  `json:"videoName"`
	Viewes    float64 `json:"viewes"`
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

func (s *service) Viewing(ctx context.Context, video string) (respose Response, err error) {
	incrementBy := 1

	res, err := s.config.Cache.Viewed(ctx, s.config.LifetimeKey(), video, float64(incrementBy))
	if err != nil {
		return respose, err
	}

	_, err = s.config.Cache.Viewed(ctx, s.config.DayKey(), video, float64(incrementBy))
	if err != nil {
		return respose, err
	}

	return Response{
		VideoName: video,
		Viewes:    res,
	}, nil
}

func (s *service) views(ctx context.Context, key string, video string) (result Response, err error) {

	views, err := s.config.Cache.ViewCount(ctx, key, video)
	if err != nil {
		return result, err
	}

	result = Response{
		VideoName: video,
		Viewes:    views,
	}

	return result, nil

}

func (s *service) DayViews(ctx context.Context, video string) (result Response, err error) {
	return s.views(ctx, s.config.DayKey(), video)
}

func (s *service) LifetimeViews(ctx context.Context, video string) (result Response, err error) {
	return s.views(ctx, s.config.LifetimeKey(), video)
}
