package leaderendpoint

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	leaderservice "github.com/SurajKadam7/leaderboard_microservice/leaderboard/service"
	"github.com/SurajKadam7/leaderboard_microservice/model"
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
	Limit string
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

		limit, err := strconv.ParseInt(l.Limit, 10, 64)
		if err != nil {
			return response, errors.New("not able to parse the limit into int")
		}

		res, err := s.DayTopViewd(ctx, limit)

		return &DayTopResponse{
			Videos: res,
			Err:    err,
		}, nil

	}
}

type LifeTimeTopRequest struct {
	Limit string
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

		limit, err := strconv.ParseInt(l.Limit, 10, 64)

		if err != nil {
			return response, errors.New("not able to parse the limit into int")
		}

		res, err := s.LifetimeTopViews(ctx, limit)

		return &LifeTimeTopResponse{
			Videos: res,
			Err:    err,
		}, nil

	}
}
