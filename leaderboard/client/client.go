package client

import (
	"context"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httpkit "github.com/go-kit/kit/transport/http"
	leaderendpoint "github.com/surajkadam/youtube_assignment/leaderboard/endpoint"
	leadertransport "github.com/surajkadam/youtube_assignment/leaderboard/transport"
	"github.com/surajkadam/youtube_assignment/model"
)

type ClientHttp struct {
	dayTopViewed, lifetimeTopViwed endpoint.Endpoint
}

// DayTopViewd this will return the day top viewed videos with the given videos limit
func (ch ClientHttp) DayTopViewd(ctx context.Context, limit string) (result []model.ViedeoDetails, err error) {

	req := leaderendpoint.DayTopRequest{
		Limit: limit,
	}
	res, err := ch.dayTopViewed(ctx, req)

	if err != nil {
		return []model.ViedeoDetails{}, err
	}

	result = res.([]model.ViedeoDetails)
	return
}

// LifetimeTopViewd this will return the Lifetime top viewed videos with the given videos limit
func (ch ClientHttp) LifetimeTopViews(ctx context.Context, limit string) (result []model.ViedeoDetails, err error) {
	req := leaderendpoint.LifeTimeTopRequest{
		Limit: limit,
	}
	res, err := ch.lifetimeTopViwed(ctx, req)

	if err != nil {
		return []model.ViedeoDetails{}, err
	}

	result = res.([]model.ViedeoDetails)
	return
}

func NewHttpClient(tgt *url.URL) ClientHttp { // DayTopViewd this will return the day top viewed videos with the given videos limit
	chttp := ClientHttp{
		dayTopViewed: httpkit.NewClient(
			"GET",
			tgt,

			leadertransport.ClientDayTopViewedEncoder(),
			leadertransport.ClientDayTopViewdDecoder(),
		).Endpoint(),
		lifetimeTopViwed: httpkit.NewClient(
			"GET",
			tgt,
			leadertransport.ClientLifetimeTopViewedEncoder(),
			leadertransport.ClientLifetimeTopViewdDecoder(),
		).Endpoint(),
	}
	return chttp
}
