package leaderendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	leaderservice "github.com/surajkadam/youtube_assignment/leaderboard/service"
	"github.com/surajkadam/youtube_assignment/model"
)

type videoViews struct {
	Name  string `json:"name,omitempty"`
	Views int64  `json:"views,omitempty"`
}

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

type DayTopRequest struct {
	Limit int64
}

type DayTopResponse struct {
	Videos []model.ViedeoDetails `json:"videos,omitempty"`
	Err    error                 `json:"error,omitempty"`
}

func (dt *DayTopResponse) Error() error {
	return dt.Err
}

func makeDayTopViewsEndpoint(s leaderservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		l := request.(*DayTopRequest)

		if err != nil {
			return nil, err
		}

		res, err := s.DayTopViewd(ctx, l.Limit)

		return &DayTopResponse{
			Videos: res,
			Err:    err,
		}, nil

	}
}

type LifeTimeTopRequest struct {
	Limit int64
}

type LifeTimeTopResponse struct {
	Videos []model.ViedeoDetails `json:"videos,omitempty"`
	Err    error                 `json:"error,omitempty"`
}

func (lt *LifeTimeTopResponse) Error() error {
	return lt.Err
}

func makeLifetimeTopViewsEndpoint(s leaderservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		l := request.(*LifeTimeTopRequest)

		if err != nil {
			return nil, err
		}
		res, err := s.LifetimeTopViews(ctx, l.Limit)

		return &LifeTimeTopResponse{
			Videos: res,
			Err:    err,
		}, nil

	}
}
