package viewendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	viewservice "github.com/surajkadam/youtube_assignment/view/service"
)

type Set struct {
	ViewingEndpoint       endpoint.Endpoint
	DayViewsEndpoint      endpoint.Endpoint
	LifetimeViewsEndpoint endpoint.Endpoint
}

func New(srv viewservice.Service) Set {
	v := makeViewingEndpoint(srv)
	d := makeDayViewsEndpoint(srv)
	l := makeLifetimeViewsEndpoint(srv)

	return Set{
		ViewingEndpoint:       v,
		DayViewsEndpoint:      d,
		LifetimeViewsEndpoint: l,
	}

}

func makeViewingEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(string)
		res, err := s.Viewing(ctx, video)

		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func makeDayViewsEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(string)
		res, err := s.DayViews(ctx, video)

		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func makeLifetimeViewsEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(string)

		res, err := s.LifetimeViews(ctx, video)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}
