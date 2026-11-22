package interfaces

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGre(t *testing.T) {
	opnsense_url := os.Getenv("OPNSENSE_URI")
	opnsense_key := os.Getenv("OPNSENSE_API_KEY")
	opnsense_secret := os.Getenv("OPNSENSE_API_SECRET")

	api_client := api.NewClient(api.Options{
		Uri:           opnsense_url,
		APIKey:        opnsense_key,
		APISecret:     opnsense_secret,
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{
		Api: api_client,
	}
	ctx := context.Background()

	// The CI VM only has the wan interface assigned.
	gre := &Gre{
		LocalAddr:        "wan",
		RemoteAddr:       "192.0.2.2",
		TunnelLocalAddr:  "10.255.255.1",
		TunnelRemoteAddr: "10.255.255.0",
		TunnelRemoteNet:  "31",
		Description:      "test gre",
	}

	key, err := controller.AddGre(ctx, gre)
	if err != nil {
		t.Fatalf("Failed to add GRE: %v", err)
	}
	t.Logf("Added GRE with key: %s", key)

	retrieved, err := controller.GetGre(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get GRE: %v", err)
	}
	t.Logf("Retrieved GRE: %+v", retrieved)
	if retrieved.LocalAddr != "wan" || retrieved.RemoteAddr != gre.RemoteAddr ||
		retrieved.TunnelLocalAddr != gre.TunnelLocalAddr || retrieved.TunnelRemoteNet != "31" {
		t.Fatalf("Retrieved GRE does not match: %+v", retrieved)
	}
	if retrieved.Device == "" {
		t.Fatalf("Expected OPNsense to assign a GRE device name")
	}

	// Local address can also be a plain IP address.
	gre.Device = retrieved.Device
	gre.LocalAddr = "192.0.2.1"
	gre.Description = "test gre updated"
	if err := controller.UpdateGre(ctx, key, gre); err != nil {
		t.Fatalf("Failed to update GRE: %v", err)
	}
	retrieved, err = controller.GetGre(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated GRE: %v", err)
	}
	if retrieved.LocalAddr != "192.0.2.1" || retrieved.Description != "test gre updated" || retrieved.Device != gre.Device {
		t.Fatalf("Updated GRE does not match: %+v", retrieved)
	}

	if err := controller.DeleteGre(ctx, key); err != nil {
		t.Fatalf("Failed to delete GRE: %v", err)
	}
	t.Logf("Deleted GRE with key: %s", key)
}
