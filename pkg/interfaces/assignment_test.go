package interfaces

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestAssignment(t *testing.T) {
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

	// The CI VM only has the wan interface assigned, so create a VLAN to assign.
	vlanKey, err := controller.AddVlan(ctx, &Vlan{
		Priority:    "0",
		Description: "Test assignment VLAN",
		Parent:      "vtnet0",
		Tag:         "3999",
	})
	if err != nil {
		t.Fatalf("Failed to add VLAN: %v", err)
	}
	defer func() {
		if err := controller.DeleteVlan(ctx, vlanKey); err != nil {
			t.Errorf("Failed to delete VLAN: %v", err)
		}
	}()
	vlan, err := controller.GetVlan(ctx, vlanKey)
	if err != nil {
		t.Fatalf("Failed to get VLAN: %v", err)
	}

	assignment := &Assignment{
		Device:              api.SelectedMap(vlan.Device),
		Description:         "testassign",
		Lock:                "0", // a locked interface cannot be deleted
		Enable:              "1",
		DisableVlanHwFilter: "0",
		Dhcp6IaPdLen:        "0",
		Dhcp6RequestDns:     "1",
		Mtu:                 "1400",
		Type4:               "staticv4",
		IpAddr:              "192.0.2.1/24",
	}

	key, err := controller.AddAssignment(ctx, assignment)
	if err != nil {
		t.Fatalf("Failed to add assignment: %v", err)
	}
	t.Logf("Added assignment with key: %s", key)

	retrieved, err := controller.GetAssignment(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get assignment: %v", err)
	}
	if retrieved.Device.String() != vlan.Device {
		t.Fatalf("Retrieved device does not match: got %s, want %s", retrieved.Device, vlan.Device)
	}
	if retrieved.Mtu != "1400" || retrieved.Type4.String() != "staticv4" || retrieved.IpAddr != "192.0.2.1/24" {
		t.Fatalf("Retrieved assignment does not match: %+v", retrieved)
	}

	assignment.Mtu = "1500"
	assignment.BlockBogons = "1"
	if err := controller.UpdateAssignment(ctx, key, assignment); err != nil {
		t.Fatalf("Failed to update assignment: %v", err)
	}
	retrieved, err = controller.GetAssignment(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated assignment: %v", err)
	}
	if retrieved.Mtu != "1500" || retrieved.BlockBogons != "1" {
		t.Fatalf("Updated assignment does not match: %+v", retrieved)
	}

	if err := controller.DeleteAssignment(ctx, key); err != nil {
		t.Fatalf("Failed to delete assignment: %v", err)
	}
	t.Logf("Deleted assignment with key: %s", key)
}
