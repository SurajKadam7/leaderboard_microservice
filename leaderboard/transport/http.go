package leadertransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	leaderendpoint "github.com/surajkadam/youtube_assignment/leaderboard/endpoint"

	"github.com/go-kit/kit/transport"
)

func NewHTTPHandler(endpoint leaderendpoint.Set, logger kitlog.Logger) http.Handler {
	// it will help loggin and error handling on the transport layer
	// it can remove the default settings and we can use our own errencoder etc
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	dayViewsHandler := kithttp.NewServer(
		endpoint.DayTopViewsEndpoint,
		decodeDayTopViewsRequest,
		encodeResponse,
		opts...,
	)

	lifetimeViewsHandler := kithttp.NewServer(
		endpoint.LifetimeTopViewsEndpoint,
		decodeLifetimeTopViewsRequest,
		encodeResponse,
		opts...,
	)

	mux := mux.NewRouter()

	mux.Handle("/top/viewes", lifetimeViewsHandler).Methods("GET")
	mux.Handle("/top/day/viewes", dayViewsHandler).Methods("GET")

	return mux

}

func decodeDayTopViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	l := r.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(l, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("invalid value in limit can not parse it: %v", limit)
	}
	return &leaderendpoint.DayTopRequest{
		Limit: limit,
	}, nil

}

func decodeLifetimeTopViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	l := r.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(l, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("invalid value in limit can not parse it: %v", limit)
	}
	return &leaderendpoint.LifeTimeTopRequest{
		Limit: limit,
	}, nil
}

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
