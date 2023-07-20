package viewendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/surajkadam/youtube_assignment/model"
	viewservice "github.com/surajkadam/youtube_assignment/view/service"
)

type Set struct {
	ViewingEndpoint       endpoint.Endpoint
	DayViewsEndpoint      endpoint.Endpoint
	LifetimeViewsEndpoint endpoint.Endpoint
	AddVideosEndpoint     endpoint.Endpoint
}

func New(srv viewservice.Service) Set {
	v := makeViewingEndpoint(srv)
	d := makeDayViewsEndpoint(srv)
	l := makeLifetimeViewsEndpoint(srv)
	a := addVideosViewsEndpoint(srv)

	return Set{
		ViewingEndpoint:       v,
		DayViewsEndpoint:      d,
		LifetimeViewsEndpoint: l,
		AddVideosEndpoint:     a,
	}
}

type videoViews struct {
	model.ViedeoDetails
	Err error `json:"error,omitempty"`
}

type ViewingRequest struct {
	Video string
}

type ViewingResponse struct {
	videoViews
}

func (vw *ViewingResponse) Error() error {
	return vw.Err
}

func makeViewingEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(ViewingRequest)
		res, err := s.View(ctx, video.Video)

		return &ViewingResponse{

			videoViews: videoViews{
				ViedeoDetails: res,
				Err:           err,
			},
		}, nil

	}
}

type DayViewsRequest struct {
	Video string
}

type DayViewsResponse struct {
	videoViews
}

func (dv *DayViewsResponse) Error() error {
	return dv.Err
}

func makeDayViewsEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(DayViewsRequest)
		res, err := s.DayViews(ctx, video.Video)

		return &DayViewsResponse{
			videoViews: videoViews{
				ViedeoDetails: res,
				Err:           err,
			},
		}, nil
	}
}

type LifetimeViewRequest struct {
	Video string
}

type LifetimeViewResponse struct {
	videoViews
}

func (lt *LifetimeViewResponse) Error() error {
	return lt.Err
}

func makeLifetimeViewsEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(LifetimeViewRequest)

		res, err := s.LifetimeViews(ctx, video.Video)

		return &LifetimeViewResponse{
			videoViews: videoViews{
				ViedeoDetails: res,
				Err:           err,
			},
		}, nil

	}
}

type AddVideosRequest struct {
	Videos []model.Video `json:"videos"`
}

type AddVideosResponse struct {
	model.AddVideoStatus
	Err  error `json:"error,omitempty"`
	Stat int32 `json:"status,omitempty"`
}

func (av *AddVideosResponse) Error() error {
	return av.Err
}

func addVideosViewsEndpoint(s viewservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		video := request.(AddVideosRequest)

		res, err := s.AddVideos(ctx, video.Videos)

		return &AddVideosResponse{
			AddVideoStatus: res,
			Err:            err,
		}, nil

	}
}
