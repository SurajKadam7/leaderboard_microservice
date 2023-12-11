// need to know how the context.WithTimeOut will work here
package integration_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/SurajKadam7/leaderboard_microservice/leaderboard/client"
	"github.com/SurajKadam7/leaderboard_microservice/model"
)

type args struct {
	limit int64
}

var tests5 = []struct {
	name               string
	args               args
	expectedStatusCode int
	wantErr            bool
	wantResult         []model.ViedeoDetails
}{
	// lifetime top views handler
	{
		name:               "top_lifetime_viewsHandler_test1",
		args:               args{limit: 10},
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult: []model.ViedeoDetails{
			{
				VideoName: "race2",
				Viewes:    2,
			},
			{
				VideoName: "race1",
				Viewes:    2,
			},
			{
				VideoName: "race3",
				Viewes:    1,
			},
		},
	},

	{
		name:               "top_lifetime_viewsHandler_test2",
		args:               args{limit: -10},
		expectedStatusCode: 200,
		wantErr:            true,
		wantResult:         make([]model.ViedeoDetails, 0),
	},

	// day top views
	{
		name:               "top_day_viewsHandler_test1",
		args:               args{limit: 10},
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult: []model.ViedeoDetails{
			{
				VideoName: "race2",
				Viewes:    2,
			},
			{
				VideoName: "race1",
				Viewes:    2,
			},
			{
				VideoName: "race3",
				Viewes:    1,
			},
		},
	},

	{
		name:               "top_day_viewsHandler_test2",
		args:               args{limit: -10},
		expectedStatusCode: 200,
		wantErr:            true,
		wantResult:         make([]model.ViedeoDetails, 0),
	},
}

func Test_LeaderBoard_Day_handlers(t *testing.T) {
	// seed data
	data := map[string]int{
		"race1": 2,
		"race2": 2,
		"race3": 1,
	}

	handler, repo, clearRedis := getHandler()
	defer clearRedis()

	for key, value := range data {
		repo.Viewed(context.Background(), key, int64(value))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()

	url, err := url.Parse(srv.URL)

	if err != nil {
		panic("error")
	}

	for _, tt := range tests5 {
		t.Run(tt.name, func(t *testing.T) {
			s := client.NewHttpClient(url)
			// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
			// defer cancel()
			limit := strconv.FormatInt(tt.args.limit, 10)

			gotResult, err := s.DayTopViewd(context.Background(), limit)
			if err != nil {
				fmt.Println("\nError is : ", err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("service.LifetimeTopViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.LifetimeTopViews() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}

}

func Test_LeaderBoard_Liftime_handlers(t *testing.T) {
	// seed data
	data := map[string]int{
		"race1": 2,
		"race2": 2,
		"race3": 1,
	}

	handler, repo, clearRedis := getHandler()
	defer clearRedis()

	for key, value := range data {
		repo.Viewed(context.Background(), key, int64(value))
	}
	// defer clearRedis()

	srv := httptest.NewServer(handler)
	defer srv.Close()

	url, err := url.Parse(srv.URL)

	if err != nil {
		panic("error")
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests5 {
		wg.Add(1)
		go t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			s := client.NewHttpClient(url)
			limit := strconv.FormatInt(tt.args.limit, 10)
			gotResult, err := s.LifetimeTopViews(context.Background(), limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.LifetimeTopViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.LifetimeTopViews() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
	wg.Wait()

}

// youtube_assignment
