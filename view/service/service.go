package viewservice

import (
	"context"

	"github.com/SurajKadam7/leaderboard_microservice/model"
	"github.com/SurajKadam7/leaderboard_microservice/repo"
)

// my bussiness logic will stay here ...
type Service interface {
	// View will increase the count of the video which the viewer is watching
	View(ctx context.Context, video string) (result model.ViedeoDetails, err error)

	//DayViews will get name of the video and return the current day views count of that videos
	DayViews(ctx context.Context, video string) (result model.ViedeoDetails, err error)

	// LifetimeViews will get name of the video and return the lifetime views count of that videos
	LifetimeViews(ctx context.Context, video string) (result model.ViedeoDetails, err error)

	// AddVideos will get the slice of videos and add them into the db. it will return status ok in success
	AddVideos(ctx context.Context, videos []model.Video) (result model.AddVideoStatus, err error)
}

type service struct {
	repo repo.Repository
}

func New(c repo.Repository) *service {
	return &service{
		repo: c,
	}
}

func (s *service) View(ctx context.Context, video string) (result model.ViedeoDetails, err error) {
	var incrementBy int64 = 1

	res, err := s.repo.Viewed(ctx, video, incrementBy)
	if err != nil {
		return result, err
	}

	return model.ViedeoDetails{
		VideoName: video,
		Viewes:    res,
	}, nil
}

func (s *service) DayViews(ctx context.Context, video string) (result model.ViedeoDetails, err error) {
	views, err := s.repo.DayViewCount(ctx, video)

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

	views, err := s.repo.LifetimeViewCount(ctx, video)

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

	err = s.repo.AddVideos(ctx, videos)
	if err != nil {
		return result, err
	}

	result = model.AddVideoStatus{
		Status: "ok",
	}
	return result, nil
}
