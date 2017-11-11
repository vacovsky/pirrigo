package main

import (
	"testing"
)

func Test_logToFile(t *testing.T) {
	SETTINGS.parseSettingsFile()
	type args struct {
		message    string
		stackTrace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				message: "teest loggyef febsifbesibfif",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logToFile(tt.args.message, tt.args.stackTrace)
		})
	}
}
