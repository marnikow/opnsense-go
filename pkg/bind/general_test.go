package bind

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestGeneral round-trips one setting and restores it. It skips when the
// os-bind plugin is not installed (the CI VM runs a stock image).
func TestGeneral(t *testing.T) {
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

	respGet, err := controller.GeneralGet(ctx)
	if err != nil || respGet.General.Port == "" {
		t.Skipf("os-bind not available: %v", err)
	}

	orig := respGet.General
	updated := orig
	updated.LogSize = "6"
	if res, err := controller.GeneralUpdate(ctx, &updated); err != nil || res.Result != "saved" {
		t.Fatalf("GeneralUpdate failed: %v %+v", err, res)
	}
	respGet, err = controller.GeneralGet(ctx)
	if err != nil {
		t.Fatalf("GeneralGet (after update) failed: %v", err)
	}
	if respGet.General.LogSize != "6" {
		t.Fatalf("logsize not updated; got %q", respGet.General.LogSize)
	}
	if respGet.General.RndcSecret != orig.RndcSecret {
		t.Fatalf("rndcsecret changed by an update that did not touch it")
	}

	if _, err := controller.GeneralUpdate(ctx, &orig); err != nil {
		t.Fatalf("GeneralUpdate (restore) failed: %v", err)
	}
}
