package hostdiscovery

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestSettings round-trips one general setting and restores it (no
// reconfigure, so the service on the test VM is untouched).
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
	if err != nil {
		t.Fatalf("SettingsGet failed: %v", err)
	}
	orig := respGet.Hostdiscovery
	if orig.General.Enabled == "" {
		t.Fatal("Expected general.enabled to be set")
	}

	updated := orig
	updated.General.Expire4Interval = "7200"
	if res, err := controller.SettingsUpdate(ctx, &updated); err != nil || res.Result != "saved" {
		t.Fatalf("SettingsUpdate failed: %v %+v", err, res)
	}
	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet (after update) failed: %v", err)
	}
	if respGet.Hostdiscovery.General.Expire4Interval != "7200" {
		t.Fatalf("Expire4Interval not updated; got %q", respGet.Hostdiscovery.General.Expire4Interval)
	}

	if _, err := controller.SettingsUpdate(ctx, &orig); err != nil {
		t.Fatalf("SettingsUpdate (restore) failed: %v", err)
	}
}
