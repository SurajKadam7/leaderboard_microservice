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
		decodeTopViewsRequest,
		encodeResponse,
		opts...,
	)

	lifetimeViewsHandler := kithttp.NewServer(
		endpoint.LifetimeTopViewsEndpoint,
		decodeTopViewsRequest,
		encodeResponse,
		opts...,
	)

	mux := mux.NewRouter()

	mux.Handle("/youtube/top/viewes", lifetimeViewsHandler)
	mux.Handle("/youtube/top/day/viewes", dayViewsHandler)

	return mux

}

func decodeTopViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	l := r.URL.Query().Get("limit")
	limit, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value in limit can not parse it: %v", limit)
	}
	return limit, nil
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

// this will return the final error to the client...
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case leaderendpoint.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})

}
