package leaderservice

import (
	"context"

	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	"github.com/surajkadam/youtube_assignment/model"
	cache "github.com/surajkadam/youtube_assignment/repository"
)

type Service interface {
	DayTopViewd(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error)
	LifetimeTopViews(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error)
}

type service struct {
	cache cache.Repository
}

func New(c cache.Repository) *service {
	return &service{
		cache: c,
	}
}

func validate(limit int64) bool {
	if limit < 0 {
		return false
	}

	if limit > 1000 {
		limit = 1000
	}

	return true
}

// need to handle key point..
func (s *service) DayTopViewd(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error) {
	ok := validate(limit)
	if !ok {
		return nil, youtubeerror.ErrInvalidLimitValue
	}
	result, err = s.cache.DayTopViewed(ctx, limit)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) LifetimeTopViews(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error) {

	ok := validate(limit)
	if !ok {
		return nil, youtubeerror.ErrInvalidLimitValue
	}

	result, err = s.cache.LifetimeTopViewed(ctx, limit)

	if err != nil {
		return nil, err
	}

	return result, nil
}
