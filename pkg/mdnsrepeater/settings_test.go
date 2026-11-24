package mdnsrepeater

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestSettings round-trips the settings and restores them (no reconfigure, and
// the repeater stays disabled). It skips when the os-mdns-repeater plugin is not
// installed (the CI VM runs a stock image).
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
	if err != nil || respGet.Mdnsrepeater.Enabled == "" {
		t.Skipf("os-mdns-repeater not available: %v", err)
	}

	orig := respGet.Mdnsrepeater
	// interfaces is required and defaults to lan, which the single-interface
	// test VM does not have; keep wan when nothing valid is selected.
	if len(orig.Interfaces) == 0 {
		orig.Interfaces = api.SelectedMapList{"wan"}
	}

	updated := orig
	updated.EnableCarp = "1"
	updated.Blocklist = api.SelectedMapList{"192.0.2.0/24", "198.51.100.7/32"}
	if res, err := controller.SettingsUpdate(ctx, &updated); err != nil || res.Result != "saved" {
		t.Fatalf("SettingsUpdate failed: %v %+v", err, res)
	}
	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet (after update) failed: %v", err)
	}
	if respGet.Mdnsrepeater.EnableCarp != "1" || len(respGet.Mdnsrepeater.Blocklist) != 2 {
		t.Fatalf("settings not updated; got %+v", respGet.Mdnsrepeater)
	}

	if _, err := controller.SettingsUpdate(ctx, &orig); err != nil {
		t.Fatalf("SettingsUpdate (restore) failed: %v", err)
	}
}
