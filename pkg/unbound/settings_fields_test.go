package unbound

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestSettingsStatsSafeSearchBelowNxDomain round-trips general.stats,
// general.safesearch and advanced.belownxdomain, then restores them.
func TestSettingsStatsSafeSearchBelowNxDomain(t *testing.T) {
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
		t.Fatalf("Failed to get settings: %v", err)
	}
	orig := respGet.Unbound

	updated := orig
	updated.General.Stats = "1"
	updated.General.SafeSearch = "1"
	updated.Advanced.BelowNxDomain = "0"
	if _, err := controller.SettingsUpdate(ctx, &updated); err != nil {
		t.Fatalf("Failed to update settings: %v", err)
	}

	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get updated settings: %v", err)
	}
	g, a := respGet.Unbound.General, respGet.Unbound.Advanced
	if g.Stats != "1" || g.SafeSearch != "1" || a.BelowNxDomain != "0" {
		t.Fatalf("settings not updated: stats %q safesearch %q belownxdomain %q", g.Stats, g.SafeSearch, a.BelowNxDomain)
	}

	if _, err := controller.SettingsUpdate(ctx, &orig); err != nil {
		t.Fatalf("Failed to restore settings: %v", err)
	}
}
