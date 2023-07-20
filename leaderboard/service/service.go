package leaderservice

import (
	"context"

	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	"github.com/surajkadam/youtube_assignment/model"
	"github.com/surajkadam/youtube_assignment/repo"
)

type Service interface {
	// DayTopViewd this will return the day top viewed videos with the given videos limit
	DayTopViewd(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error)

	// LifetimeTopViewd this will return the Lifetime top viewed videos with the given videos limit
	LifetimeTopViews(ctx context.Context, limit int64) (result []model.ViedeoDetails, err error)
}

type service struct {
	repo repo.Repository
}

func New(c repo.Repository) *service {
	return &service{
		repo: c,
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
	result, err = s.repo.DayTopViewed(ctx, limit)

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

	result, err = s.repo.LifetimeTopViewed(ctx, limit)

	if err != nil {
		return nil, err
	}

	return result, nil
}
