package dyndns

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestSettings round-trips daemon_delay and restores it. It skips when the
// os-ddclient plugin is not installed (the CI VM runs a stock image).
func TestSettings(t *testing.T) {
	controller := Controller{Api: api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})}
	ctx := context.Background()

	respGet, err := controller.SettingsGet(ctx)
	if err != nil || respGet.Ddclient.General.DaemonDelay == "" {
		t.Skipf("os-ddclient not available: %v", err)
	}
	orig := respGet.Ddclient

	updated := orig
	updated.General.DaemonDelay = "301"
	if res, err := controller.SettingsUpdate(ctx, &updated); err != nil || res.Result != "saved" {
		t.Fatalf("SettingsUpdate failed: %v %+v", err, res)
	}
	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet (after update) failed: %v", err)
	}
	if respGet.Ddclient.General.DaemonDelay != "301" {
		t.Fatalf("daemon_delay not updated; got %q", respGet.Ddclient.General.DaemonDelay)
	}

	if _, err := controller.SettingsUpdate(ctx, &orig); err != nil {
		t.Fatalf("SettingsUpdate (restore) failed: %v", err)
	}
}
