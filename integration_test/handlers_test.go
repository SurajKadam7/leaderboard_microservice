package integration_test

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var tests = []struct {
	name               string
	url                string
	method             string
	payload            string
	expectedStatusCode int
	wantErr            bool
	wantResult         map[string]any
}{
	// viewed handler ...
	{
		name:               "viwedHandler_test1",
		url:                "/video/viewing?name=race1",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race1", "viewes": 1.0},
	},
	{
		name:               "viwedHandler_test2",
		url:                "/video/viewing?name=race1",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race1", "viewes": 2.0},
	},
	{
		name:               "viwedHandler_test3",
		url:                "/video/viewing?name=race2",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race2", "viewes": 1.0},
	},
	{
		name:               "viwedHandler_test4",
		url:                "/video/viewing?name=race2",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race2", "viewes": 2.0},
	},
	{
		name:               "viwedHandler_test5",
		url:                "/video/viewing?name=race3",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race3", "viewes": 1.0},
	},

	// lifetime views handler
	{
		name:               "lifetime_viewsHandler_test1",
		url:                "/video/viewes?name=race1",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race1", "viewes": 2.0},
	},
	{
		name:               "lifetime_viewsHandler_test2",
		url:                "/video/viewes?name=race2",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race2", "viewes": 2.0},
	},
	{
		name:               "lifetime_viewsHandler_test3",
		url:                "/video/viewes?name=race3",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race3", "viewes": 1.0},
	},

	// day views handler
	{
		name:               "day_viewsHandler_test1",
		url:                "/video/day/viewes?name=race1",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race1", "viewes": 2.0},
	},
	{
		name:               "day_viewsHandler_test2",
		url:                "/video/day/viewes?name=race2",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race2", "viewes": 2.0},
	},
	{
		name:               "day_viewsHandler_test3",
		url:                "/video/day/viewes?name=race3",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"videoName": "race3", "viewes": 1.0},
	},
}

// here I have not used the clearRedis call back function just to keep data in the redis
// after running the test 2 which is dependent on the redis data I have called
// the call back function into it

// sequncy of running these testCases is always 1, 2, 3

func Test_Views_handlers(t *testing.T) {
	handler, _ := getHandler()
	// defer clearRedis()

	srv := httptest.NewTLSServer(handler)
	defer srv.Close()

	for _, test := range tests {
		url := srv.URL + test.url

		resp, err := srv.Client().Get(url)

		if (err != nil) && !test.wantErr {
			t.Errorf("testName : %s, error = %v, wantErr %v", test.name, err, test.wantErr)
			continue
		}

		gotResult := map[string]any{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResult)

		if !reflect.DeepEqual(gotResult, test.wantResult) {
			t.Errorf("testName %v = %v, want %v", test.name, gotResult, test.wantResult)
		}

		resp.Body.Close()

	}

}

var tests2 = []struct {
	name               string
	url                string
	method             string
	payload            string
	expectedStatusCode int
	wantErr            bool
	wantResult         []map[string]any
}{
	// lifetime top views handler
	{
		name:               "top_lifetime_viewsHandler_test1",
		url:                "/top/viewes?limit=100000",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         []map[string]any{{"videoName": "race2", "viewes": 2.0}, {"videoName": "race1", "viewes": 2.0}, {"videoName": "race3", "viewes": 1.0}},
	},

	{
		name:               "top_lifetime_viewsHandler_test2",
		url:                "/top/viewes?limit=2",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         []map[string]any{{"videoName": "race2", "viewes": 2.0}, {"videoName": "race1", "viewes": 2.0}},
	},

	// day top views
	{
		name:               "top_day_viewsHandler_test1",
		url:                "/top/day/viewes?limit=100",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         []map[string]any{{"videoName": "race2", "viewes": 2.0}, {"videoName": "race1", "viewes": 2.0}, {"videoName": "race3", "viewes": 1.0}},
	},

	{
		name:               "top_day_viewsHandler_test2",
		url:                "/top/day/viewes?limit=2",
		method:             "GET",
		payload:            "",
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         []map[string]any{{"videoName": "race2", "viewes": 2.0}, {"videoName": "race1", "viewes": 2.0}},
	},
}

func Test_LeaderBoard_handlers(t *testing.T) {
	handler, clearRedis := getHandler()
	defer clearRedis()

	srv := httptest.NewTLSServer(handler)
	defer srv.Close()

	for _, test := range tests2 {
		url := srv.URL + test.url

		resp, err := srv.Client().Get(url)

		if (err != nil) && !test.wantErr {
			t.Errorf("testName : %s, error = %v, wantErr %v", test.name, err, test.wantErr)
			continue
		}

		// handled validation for both post and get

		var result map[string][]map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)
		gotResult := result["videos"]

		if !reflect.DeepEqual(gotResult, test.wantResult) {
			t.Errorf("testName %v = %v, want %v", test.name, gotResult, test.wantResult)
		}

		resp.Body.Close()

	}

}

var tests3 = []struct {
	name               string
	url                string
	method             string
	payload            string
	expectedStatusCode int
	wantErr            bool
	wantResult         map[string]any
}{
	// add Videos handler localhost:8080/video/add
	{
		name:   "add_videos_Handler_test1",
		url:    "/video/add",
		method: "POST",
		payload: `{"videos": [
									{"name": "Movie1"},
									{"name": "Movie2"},
									{"name": "Movie3"}
									]
						}`,
		expectedStatusCode: 200,
		wantErr:            false,
		wantResult:         map[string]any{"Status": "ok"},
	},
	// syntax error in the json payload not started array bracket.
	{
		name:   "add_videos_Handler_test2",
		url:    "/video/add",
		method: "POST",
		payload: `{"videos": 
									{"name1": "Movie1"},
									{"name": "Movie2"},
									{"name": "Movie3"}
									]
						}`,
		expectedStatusCode: 200,
		wantErr:            true,
		wantResult:         map[string]any{"error": "invalid json payload, not able to parse"},
	},
}

func Test_Add_handlers(t *testing.T) {
	handler, clearRedis := getHandler()
	defer clearRedis()

	srv := httptest.NewTLSServer(handler)
	defer srv.Close()

	for _, test := range tests3 {
		url := srv.URL + test.url
		reder := strings.NewReader(test.payload)
		resp, _ := srv.Client().Post(url, "application/json", reder)

		gotResult := map[string]any{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&gotResult)

		err, ok := gotResult["error"]

		if !test.wantErr && ok {
			t.Errorf("testName : %s, error = %v, wantErr %v", test.name, err, test.wantErr)
			continue
		}

		if !reflect.DeepEqual(gotResult, test.wantResult) {
			t.Errorf("testName %v = %v, want %v", test.name, gotResult, test.wantResult)
		}

		resp.Body.Close()

	}

}
