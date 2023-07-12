package viewtransport

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	viewendpoint "github.com/surajkadam/youtube_assignment/view/endpoint"

	"github.com/go-kit/kit/transport"
)

func NewHTTPHandler(endpoint viewendpoint.Set, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	viewingHandler := kithttp.NewServer(
		endpoint.ViewingEndpoint,
		decodeViewsRequest,
		encodeResponse,
		opts...,
	)

	dayViewsHandler := kithttp.NewServer(
		endpoint.DayViewsEndpoint,
		decodeViewsRequest,
		encodeResponse,
		opts...,
	)

	lifetimeViewsHandler := kithttp.NewServer(
		endpoint.LifetimeViewsEndpoint,
		decodeViewsRequest,
		encodeResponse,
		opts...,
	)

	mux := mux.NewRouter()

	mux.Handle("/youtube/video/viewing", viewingHandler)
	mux.Handle("/youtube/video/viewes", lifetimeViewsHandler)
	mux.Handle("/youtube/video/day/viewes", dayViewsHandler)

	return mux

}

func decodeViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	video := r.URL.Query().Get("name")
	if len(video) == 0 {
		return nil, youtubeerror.ErrEmptyVideoValuePassed
	}
	return video, nil
}

// incomplete implimentation
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
