package main

import (
	"testing"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name          string
		wantConsulUrl string
		wantErr       bool
	}{
		{
			name:          "config_testing",
			wantConsulUrl: "http://consul:8500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConsulUrl, err := loadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConsulUrl != tt.wantConsulUrl {
				t.Errorf("loadConfig() = %v, want %v", gotConsulUrl, tt.wantConsulUrl)
			}
		})
	}
}
