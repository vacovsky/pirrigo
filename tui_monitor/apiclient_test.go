package main

import (
	"testing"

	"github.com/vacovsky/pirrigo/pirri"
)

// func Test_makeGetCall(t *testing.T) {
// 	type args struct {
// 		url    string
// 		target interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    pirri.RunStatus
// 		wantErr bool
// 	}{
// 		{
// 			name: "call works",
// 			args: args{
// 				url:    "http://192.168.111.130/status/run",
// 				target: pirri.RunStatus{},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := makeGetCall(tt.args.url, tt.args.target)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("makeGetCall() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("makeGetCall() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_calcTimeDiff(t *testing.T) {
	type args struct {
		status pirri.RunStatus
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := calcTimeDiff(tt.args.status)
			if got != tt.want {
				t.Errorf("calcTimeDiff() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("calcTimeDiff() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
