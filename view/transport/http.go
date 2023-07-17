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
		decodeViewingRequest,
		encodeResponse,
		opts...,
	)

	dayViewsHandler := kithttp.NewServer(
		endpoint.DayViewsEndpoint,
		decodeDayViewRequest,
		encodeResponse,
		opts...,
	)

	lifetimeViewsHandler := kithttp.NewServer(
		endpoint.LifetimeViewsEndpoint,
		decodeLifetimeViewsRequest,
		encodeResponse,
		opts...,
	)

	addVideoHandler := kithttp.NewServer(
		endpoint.AddVideosEndpoint,
		decodeAddViedosRequest,
		encodeResponse,
		opts...,
	)
	mux := mux.NewRouter()

	mux.Handle("/video/viewing", viewingHandler).Methods("GET")
	mux.Handle("/video/viewes", lifetimeViewsHandler).Methods("GET")
	mux.Handle("/video/day/viewes", dayViewsHandler).Methods("GET")
	mux.Handle("/video/add", addVideoHandler).Methods("POST")

	return mux

}

func viewDecoder(r *http.Request) (interface{}, error) {
	video := r.URL.Query().Get("name")
	if len(video) == 0 {
		return nil, youtubeerror.ErrEmptyVideoValuePassed
	}
	return video, nil
}

func decodeViewingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	v, err := viewDecoder(r)
	if err != nil {
		return nil, err
	}
	video := v.(string)
	return viewendpoint.ViewingRequest{
		Video: video,
	}, nil
}

func decodeDayViewRequest(_ context.Context, r *http.Request) (interface{}, error) {
	v, err := viewDecoder(r)
	if err != nil {
		return nil, err
	}
	video := v.(string)
	return viewendpoint.DayViewsRequest{
		Video: video,
	}, nil
}

func decodeLifetimeViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	v, err := viewDecoder(r)
	if err != nil {
		return nil, err
	}
	video := v.(string)
	return viewendpoint.LifetimeViewRequest{
		Video: video,
	}, nil
}

func decodeAddViedosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	videos := viewendpoint.AddVideosRequest{}

	if err := json.NewDecoder(r.Body).Decode(&videos); err != nil {
		return nil, youtubeerror.ErrNotAbleToParse
	}
	return videos, nil
}

// incomplete implimentation
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Errorr); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type Errorr interface {
	Error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case youtubeerror.ErrEmptyVideoValuePassed, youtubeerror.ErrInvalidLimitValue, youtubeerror.ErrVideoNotFound:
		w.WriteHeader(http.StatusBadRequest)

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
