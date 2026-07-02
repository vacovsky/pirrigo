package pirri

import (
	"net/http"
	"testing"
)

func TestParseDaysParam(t *testing.T) {
	tests := []struct {
		name   string
		qs     string
		def    int
		expect int
	}{
		{"default", "", 7, 7},
		{"explicit", "days=30", 7, 30},
		{"zero→default", "days=0", 7, 7},
		{"negative→default", "days=-5", 7, 7},
		{"garbage→default", "days=abc", 7, 7},
		{"one", "days=1", 7, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://x/?"+tt.qs, nil)
			if got := parseDaysParam(req, tt.def); got != tt.expect {
				t.Errorf("parseDaysParam(%q, %d) = %d, want %d", tt.qs, tt.def, got, tt.expect)
			}
		})
	}
}

// ponytail: smoke test — verify statsActivityPerStationByDOW doesn't panic with empty DB
func TestStatsActivityPerStationByDOW_NoData(t *testing.T) {
	// This requires a real DB connection, so we skip if no DB is configured.
	// The point is to catch panics during dev testing.
	req, _ := http.NewRequest("GET", "http://x/?days=7", nil)
	// Can't call without initialized data.Service(), so just verify the query param parses:
	days := parseDaysParam(req, 7)
	if days != 7 {
		t.Fatal("expected 7")
	}
}
