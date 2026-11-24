package radvd

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestEntry adds, updates, finds and deletes a router advertisement entry on wan.
// The entry stays disabled, so radvd never sends advertisements on the test
// VM's network.
func TestEntry(t *testing.T) {
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

	entry := &Entry{
		Enabled:     "0",
		Interface:   api.SelectedMap("wan"),
		Mode:        api.SelectedMap("unmanaged"),
		DNS:         "1",
		MinInterval: "200",
		MaxInterval: "600",
		Priority:    api.SelectedMap("medium"),
		HopLimit:    "64",
	}

	id, err := controller.AddEntry(ctx, entry)
	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}
	t.Cleanup(func() {
		if err := controller.DeleteEntry(ctx, id); err != nil {
			t.Errorf("DeleteEntry failed: %v", err)
		}
	})

	entry.Mode = api.SelectedMap("assist")
	entry.Priority = api.SelectedMap("high")
	entry.RDNSS = api.SelectedMapList{"2001:db8::53"}
	entry.DNSSL = api.SelectedMapList{"example.com"}
	entry.Routes = api.SelectedMapList{"2001:db8:1::/64"}
	entry.LinkMTU = "1400"
	if err := controller.UpdateEntry(ctx, id, entry); err != nil {
		t.Fatalf("UpdateEntry failed: %v", err)
	}

	got, err := controller.GetEntry(ctx, id)
	if err != nil {
		t.Fatalf("GetEntry failed: %v", err)
	}
	if got.Interface != "wan" || got.Mode != "assist" || got.Priority != "high" || got.LinkMTU != "1400" ||
		len(got.RDNSS) != 1 || got.RDNSS[0] != "2001:db8::53" || len(got.DNSSL) != 1 || len(got.Routes) != 1 {
		t.Fatalf("entry not updated; got %+v", got)
	}

	search, err := controller.SearchEntry(ctx, "-1")
	if err != nil {
		t.Fatalf("SearchEntry failed: %v", err)
	}
	found := false
	for _, row := range search.Rows {
		if row.Id == id {
			found = true
		}
	}
	if !found {
		t.Fatalf("entry %s not found by SearchEntry: %+v", id, search.Rows)
	}
}
