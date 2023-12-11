package viewservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	youtubeerror "github.com/SurajKadam7/leaderboard_microservice/errors"
	"github.com/SurajKadam7/leaderboard_microservice/model"
	cache "github.com/SurajKadam7/leaderboard_microservice/repo"
	mock_cache "github.com/SurajKadam7/leaderboard_microservice/repo/mock"
)

var videos = []model.ViedeoDetails{
	{
		VideoName: "race1",
		Viewes:    10,
	},
	{
		VideoName: "race2",
		Viewes:    100,
	},
	{
		VideoName: "race3",
		Viewes:    1000,
	},
	{
		VideoName: "race4",
		Viewes:    10000,
	},
}

func Test_service_Viewing(t *testing.T) {
	// mocking layer ...
	ctr := gomock.NewController(t)
	// no need to call this ctr.Finish from go 1.14+
	defer ctr.Finish()

	cacheMock := mock_cache.NewMockRepository(ctr)
	// mock 1
	cacheMock.EXPECT().Viewed(gomock.Any(), "race3", int64(1)).Times(1).Return(1000.0, nil)
	// mock 2
	cacheMock.EXPECT().Viewed(gomock.Any(), "race4", int64(1)).Times(1).Return(0.0, youtubeerror.ErrEmptyVideoValuePassed)

	type fields struct {
		cache cache.Repository
	}
	type args struct {
		ctx   context.Context
		video string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult model.ViedeoDetails
		wantErr    bool
	}{
		{
			name:   "basic",
			fields: fields{cache: cacheMock},
			args: args{
				ctx:   context.TODO(),
				video: "race3",
			},
			wantResult: videos[2],
			wantErr:    false,
		},
		{
			name:   "empty video error",
			fields: fields{cache: cacheMock},
			args: args{
				ctx:   context.TODO(),
				video: "race4",
			},
			wantResult: model.ViedeoDetails{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.cache,
			}
			gotResult, err := s.View(tt.args.ctx, tt.args.video)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Viewing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.Viewing() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_service_DayViews(t *testing.T) {
	// mocking Repository
	ctr := gomock.NewController(t)

	mockCache := mock_cache.NewMockRepository(ctr)

	// mock 1
	mockCache.EXPECT().DayViewCount(gomock.Any(), "race1").Times(1).Return(10.0, nil)
	// mock 2
	mockCache.EXPECT().DayViewCount(gomock.Any(), "race5").Times(1).Return(0.0, youtubeerror.ErrVideoNotFound)

	type fields struct {
		cache cache.Repository
	}
	type args struct {
		ctx   context.Context
		video string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult model.ViedeoDetails
		wantErr    bool
	}{
		{
			name:       "basic",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), video: "race1"},
			wantResult: videos[0],
			wantErr:    false,
		},
		{
			name:       "basic",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), video: "race5"},
			wantResult: model.ViedeoDetails{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.cache,
			}
			gotResult, err := s.DayViews(tt.args.ctx, tt.args.video)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.DayViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.DayViews() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_service_LifetimeViews(t *testing.T) {

	mockCache := mock_cache.NewMockRepository(gomock.NewController(t))

	// mock 1
	mockCache.EXPECT().LifetimeViewCount(context.TODO(), "race1").Times(1).Return(10.0, nil)
	// mock 2
	mockCache.EXPECT().LifetimeViewCount(context.TODO(), "race5").Times(1).Return(0.0, youtubeerror.ErrVideoNotFound)

	type fields struct {
		cache cache.Repository
	}
	type args struct {
		ctx   context.Context
		video string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult model.ViedeoDetails
		wantErr    bool
	}{
		{
			name:       "basic",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), video: "race1"},
			wantResult: videos[0],
			wantErr:    false,
		},
		{
			name:       "basic",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), video: "race5"},
			wantResult: model.ViedeoDetails{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.cache,
			}
			gotResult, err := s.LifetimeViews(tt.args.ctx, tt.args.video)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.LifetimeViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.LifetimeViews() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_service_AddVideos(t *testing.T) {
	mockCache := mock_cache.NewMockRepository(gomock.NewController(t))

	arguments := [][]model.Video{
		{
			{Name: "race1"},
			{Name: "race2"},
			{Name: "race3"},
		},
	}

	// mock 1
	mockCache.EXPECT().AddVideos(
		context.TODO(),
		arguments[0],
	).Times(1).Return(nil)

	// mock 2
	mockCache.EXPECT().AddVideos(
		context.TODO(),
		nil,
	).Times(1).Return(errors.New("passed null value"))

	type fields struct {
		cache cache.Repository
	}
	type args struct {
		ctx    context.Context
		videos []model.Video
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult model.AddVideoStatus
		wantErr    bool
	}{
		{
			name:   "basic",
			fields: fields{cache: mockCache},
			args: args{
				ctx:    context.TODO(),
				videos: arguments[0],
			},
			wantResult: model.AddVideoStatus{Status: "ok"},
			wantErr:    false,
		},

		{
			name:   "error",
			fields: fields{cache: mockCache},
			args: args{
				ctx:    context.TODO(),
				videos: nil,
			},
			wantResult: model.AddVideoStatus{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.cache,
			}
			gotResult, err := s.AddVideos(tt.args.ctx, tt.args.videos)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.AddVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.AddVideos() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
