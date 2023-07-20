package leaderservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/surajkadam/youtube_assignment/model"
	cache "github.com/surajkadam/youtube_assignment/repo"
	mock_cache "github.com/surajkadam/youtube_assignment/repo/mock"
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

func Test_service_DayTopViewd(t *testing.T) {
	v := []model.ViedeoDetails{}
	for i := range videos {
		v = append(v, videos[i])
	}

	mockCache := mock_cache.NewMockRepository(gomock.NewController(t))
	// mock 1
	mockCache.EXPECT().DayTopViewed(context.TODO(), int64(5)).Times(1).Return(v, nil)
	// mock 2
	mockCache.EXPECT().DayTopViewed(context.TODO(), int64(10)).Times(1).Return(*new([]model.ViedeoDetails), errors.New("database down"))

	type fields struct {
		cache cache.Repository
	}
	type args struct {
		ctx   context.Context
		limit int64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []model.ViedeoDetails
		wantErr    bool
	}{
		{
			name:       "basic",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), limit: 5},
			wantResult: videos,
			wantErr:    false,
		},
		{
			name:       "error",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), limit: 10},
			wantResult: *new([]model.ViedeoDetails),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.cache,
			}
			gotResult, err := s.DayTopViewd(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.DayTopViewd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("service.DayTopViewd() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_service_LifetimeTopViews(t *testing.T) {
	v := []model.ViedeoDetails{}
	for i := range videos {
		v = append(v, videos[i])
	}

	mockCache := mock_cache.NewMockRepository(gomock.NewController(t))
	// mock 1
	mockCache.EXPECT().LifetimeTopViewed(context.TODO(), int64(5)).Times(1).Return(v, nil)
	// mock 2
	mockCache.EXPECT().LifetimeTopViewed(context.TODO(), int64(10)).Times(1).Return(*new([]model.ViedeoDetails), errors.New("database down"))

	type fields struct {
		cache cache.Repository
	}
	type args struct {
		ctx   context.Context
		limit int64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []model.ViedeoDetails
		wantErr    bool
	}{
		{
			name:       "basic",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), limit: 5},
			wantResult: videos,
			wantErr:    false,
		},
		{
			name:       "error",
			fields:     fields{cache: mockCache},
			args:       args{ctx: context.TODO(), limit: 10},
			wantResult: *new([]model.ViedeoDetails),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.fields.cache,
			}
			gotResult, err := s.LifetimeTopViews(tt.args.ctx, tt.args.limit)
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
