package leadertransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"reflect"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	youtubeerror "github.com/surajkadam/youtube_assignment/errors"
	leaderendpoint "github.com/surajkadam/youtube_assignment/leaderboard/endpoint"
	"github.com/surajkadam/youtube_assignment/model"

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
	limit := r.URL.Query().Get("limit")

	return &leaderendpoint.DayTopRequest{
		Limit: limit,
	}, nil

}

func decodeLifetimeTopViewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	limit := r.URL.Query().Get("limit")

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
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

// this will set the body in case of the post request
// we can add the path in the url as well here
func ClientDayTopViewedEncoder() kithttp.EncodeRequestFunc {
	return func(ctx context.Context, r *http.Request, i interface{}) error {
		r.URL.Path = path.Join(r.URL.Path, "/top/viewes")
		return createQueryFromStruct(r, i)

	}
}

type Response struct {
	Videos []model.ViedeoDetails `json:"videos,omitempty"`
	Err    string                `json:"error,omitempty"`
}

// this will decode the request
func ClientDayTopViewdDecoder() kithttp.DecodeResponseFunc {
	return func(ctx context.Context, r *http.Response) (response interface{}, err error) {
		var resStruct Response
		err = json.NewDecoder(r.Body).Decode(&resStruct)

		if err != nil {
			return
		}
		if resStruct.Err != "" {
			err = errors.New(resStruct.Err)
		}

		return resStruct.Videos, err
	}
}

func ClientLifetimeTopViewedEncoder() kithttp.EncodeRequestFunc {
	return func(ctx context.Context, r *http.Request, i interface{}) error {
		r.URL.Path = path.Join(r.URL.Path, "/top/day/viewes")
		return createQueryFromStruct(r, i)
	}
}

// this will decode the request
func ClientLifetimeTopViewdDecoder() kithttp.DecodeResponseFunc {
	return func(ctx context.Context, r *http.Response) (response interface{}, err error) {
		var resStruct Response
		err = json.NewDecoder(r.Body).Decode(&resStruct)
		if err != nil {
			return
		}
		if resStruct.Err != "" {
			err = errors.New(resStruct.Err)
		}
		return resStruct.Videos, err
	}
}

func createQueryFromStruct(r *http.Request, i interface{}) error {
	req := i

	value := reflect.ValueOf(req)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldName = strings.ToLower(fieldName)

		filedValue := value.Field(i).Interface().(string)

		queryMap := r.URL.Query()
		queryMap.Add(fieldName, filedValue)
		r.URL.RawQuery = queryMap.Encode()
	}
	return nil
}
