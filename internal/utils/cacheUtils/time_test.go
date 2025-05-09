package cacheUtils

import (
	"testing"
	"time"
)

func TestSeconds(t *testing.T) {
	type args struct {
		duration     int
		durationType time.Duration
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"test1", args{duration: 1, durationType: time.Second}, 1},
		{"test2", args{duration: 1, durationType: time.Minute}, 60},
		{"test3", args{duration: 5, durationType: time.Minute}, 300},
		{"test4", args{duration: 24, durationType: time.Hour}, 86400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Seconds(tt.args.duration, tt.args.durationType); got != tt.want {
				t.Errorf("Seconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
