package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSourceNatMode(t *testing.T) {
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

	respGet, err := controller.SourceNatGet(ctx)
	if err != nil {
		t.Fatalf("SourceNatGet failed: %v", err)
	}
	orig := respGet.Filter
	if orig.General.SnatMode == "" {
		t.Fatal("Expected general.snat_mode to be set")
	}

	updated := orig
	updated.General.SnatMode = api.SelectedMap("hybrid")
	if res, err := controller.SourceNatUpdate(ctx, &updated); err != nil || res.Result != "saved" {
		t.Fatalf("SourceNatUpdate failed: %v %+v", err, res)
	}
	respGet, err = controller.SourceNatGet(ctx)
	if err != nil {
		t.Fatalf("SourceNatGet (after update) failed: %v", err)
	}
	if respGet.Filter.General.SnatMode != "hybrid" {
		t.Fatalf("snat_mode not updated; got %q", respGet.Filter.General.SnatMode)
	}

	if _, err := controller.SourceNatUpdate(ctx, &orig); err != nil {
		t.Fatalf("SourceNatUpdate (restore) failed: %v", err)
	}
	if _, err := controller.SourceNatApply(ctx); err != nil {
		t.Fatalf("SourceNatApply failed: %v", err)
	}
}
