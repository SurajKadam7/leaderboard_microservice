package leaderendpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	leaderservice "github.com/surajkadam/youtube_assignment/leaderboard/service"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Set struct {
	DayTopViewsEndpoint      endpoint.Endpoint
	LifetimeTopViewsEndpoint endpoint.Endpoint
}

func New(srv leaderservice.Service) Set {

	d := makeDayTopViewsEndpoint(srv)
	l := makeLifetimeTopViewsEndpoint(srv)

	return Set{
		DayTopViewsEndpoint:      d,
		LifetimeTopViewsEndpoint: l,
	}

}

func makeDayTopViewsEndpoint(s leaderservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		limit := request.(int64)
		if err != nil {
			return nil, err
		}

		res, err := s.DayTopViewd(ctx, limit)

		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func makeLifetimeTopViewsEndpoint(s leaderservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		limit := request.(int64)

		if err != nil {
			return nil, err
		}
		res, err := s.LifetimeTopViews(ctx, limit)

		if err != nil {
			return nil, err
		}

		return res, nil
	}
}
