package viewservice

import (
	"context"

	"github.com/surajkadam/youtube_assignment/model"
	cache "github.com/surajkadam/youtube_assignment/repository"
)

// my bussiness logic will stay here ...
type Service interface {
	Viewing(ctx context.Context, video string) (result model.ViedeoDetails, err error)
	DayViews(ctx context.Context, video string) (result model.ViedeoDetails, err error)
	LifetimeViews(ctx context.Context, video string) (result model.ViedeoDetails, err error)
	AddVideos(ctx context.Context, videos []model.Video) (result model.AddVideoStatus, err error)
}

type service struct {
	cache cache.Repository
}

func New(c cache.Repository) *service {
	return &service{
		cache: c,
	}
}

func (s *service) Viewing(ctx context.Context, video string) (result model.ViedeoDetails, err error) {
	var incrementBy int64 = 1

	res, err := s.cache.Viewed(ctx, video, incrementBy)
	if err != nil {
		return result, err
	}

	return model.ViedeoDetails{
		VideoName: video,
		Viewes:    res,
	}, nil
}

func (s *service) DayViews(ctx context.Context, video string) (result model.ViedeoDetails, err error) {
	views, err := s.cache.DayViewCount(ctx, video)

	if err != nil {
		return result, err
	}

	result = model.ViedeoDetails{
		VideoName: video,
		Viewes:    views,
	}

	return result, nil
}

func (s *service) LifetimeViews(ctx context.Context, video string) (result model.ViedeoDetails, err error) {

	views, err := s.cache.LifetimeViewCount(ctx, video)

	if err != nil {
		return result, err
	}

	result = model.ViedeoDetails{
		VideoName: video,
		Viewes:    views,
	}

	return result, nil
}

func (s *service) AddVideos(ctx context.Context, videos []model.Video) (result model.AddVideoStatus, err error) {

	err = s.cache.AddVideos(ctx, videos)
	if err != nil {
		return result, err
	}

	result = model.AddVideoStatus{
		Status: "ok",
	}
	return result, nil
}
