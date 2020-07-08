package logging

import (
	"reflect"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestPirriLogger_tailLogs(t *testing.T) {
	type fields struct {
		lock   sync.Mutex
		logger *zap.Logger
	}
	type args struct {
		lines int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PirriLogger{
				lock:   tt.fields.lock,
				logger: tt.fields.logger,
			}
			got, err := l.tailLogs(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("PirriLogger.tailLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PirriLogger.tailLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}
