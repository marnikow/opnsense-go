package interfaces

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

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
	t.Logf("SettingsGet (initial): %+v", respGet)
	if respGet.Settings.Dhcp6RaTimeout == "" || respGet.Settings.DisableVlanHwFilter == "" {
		t.Fatal("Expected dhcp6_ratimeout and disablevlanhwfilter to be set")
	}

	// Only the stored configuration is changed; Reconfigure (a networking
	// reload) is not called.
	orig := respGet.Settings
	updated := orig
	updated.Dhcp6RaTimeout = "11"
	respSet, err := controller.SettingsUpdate(ctx, &updated)
	if err != nil {
		t.Fatalf("SettingsUpdate failed: %v", err)
	}
	if respSet.Result != "saved" {
		t.Fatalf("Expected result 'saved', got %q", respSet.Result)
	}

	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet (after update) failed: %v", err)
	}
	if respGet.Settings.Dhcp6RaTimeout != "11" {
		t.Fatalf("dhcp6_ratimeout not updated; got %q, want %q", respGet.Settings.Dhcp6RaTimeout, "11")
	}

	if _, err := controller.SettingsUpdate(ctx, &orig); err != nil {
		t.Fatalf("SettingsUpdate (restore) failed: %v", err)
	}
	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet (after restore) failed: %v", err)
	}
	if respGet.Settings.Dhcp6RaTimeout != orig.Dhcp6RaTimeout {
		t.Fatalf("dhcp6_ratimeout not restored; got %q, want %q", respGet.Settings.Dhcp6RaTimeout, orig.Dhcp6RaTimeout)
	}
}
