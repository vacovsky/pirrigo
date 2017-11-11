package main

import "testing"

func Test_testLogging(t *testing.T) {
	SETTINGS.parseSettingsFile()
	tests := []struct {
		name string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogging()
		})
	}
}
