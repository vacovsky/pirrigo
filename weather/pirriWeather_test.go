package weather

import (
	"testing"
)

func TestCurrentWeather_NoCoords(t *testing.T) {
	// Settings singleton has zero-value coords (no config file loaded),
	// so Current() should short-circuit with Status == "Error".
	resp := Service().Current()
	if resp.Status != "Error" {
		t.Errorf("expected Status=Error with no coords, got %q", resp.Status)
	}
	if resp.Temperature != 0 {
		t.Errorf("expected Temperature=0, got %f", resp.Temperature)
	}
}
