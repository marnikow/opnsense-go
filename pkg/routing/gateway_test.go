package routing

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGateway(t *testing.T) {
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

	// The CI VM only has wan; use a far gateway so the address does not need
	// to be inside the wan subnet, and keep monitoring off.
	gateway := &Gateway{
		Disabled:       "0",
		Name:           "TEST_GW",
		Description:    "test gateway",
		Interface:      "wan",
		IpProtocol:     "inet",
		Gateway:        "192.0.2.254",
		FarGateway:     "1",
		MonitorDisable: "1",
		Priority:       "200",
		Weight:         "1",
	}
	key, err := controller.AddGateway(ctx, gateway)
	if err != nil {
		t.Fatalf("Failed to add gateway: %v", err)
	}
	t.Logf("Added gateway with key: %s", key)
	defer func() {
		if err := controller.DeleteGateway(ctx, key); err != nil {
			t.Errorf("Failed to delete gateway: %v", err)
		}
	}()

	got, err := controller.GetGateway(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get gateway: %v", err)
	}
	if got.Name != "TEST_GW" || got.Interface != "wan" || got.Gateway != "192.0.2.254" || got.Priority != "200" || got.FarGateway != "1" {
		t.Fatalf("Retrieved gateway does not match: %+v", got)
	}

	search, err := controller.GatewaysSearch(ctx)
	if err != nil {
		t.Fatalf("Failed to search gateways: %v", err)
	}
	found := false
	for _, row := range search.Rows {
		if row.Uuid == key && row.Interface == "wan" {
			found = true
		}
	}
	if !found {
		t.Fatalf("gateway %s not found with interface wan in search rows", key)
	}

	gateway.Weight = "3"
	gateway.Description = "test gateway updated"
	if err := controller.UpdateGateway(ctx, key, gateway); err != nil {
		t.Fatalf("Failed to update gateway: %v", err)
	}
	got, err = controller.GetGateway(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated gateway: %v", err)
	}
	if got.Weight != "3" || got.Description != "test gateway updated" {
		t.Fatalf("Updated gateway does not match: %+v", got)
	}
}
