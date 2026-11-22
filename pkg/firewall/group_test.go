package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGroup(t *testing.T) {
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

	// The CI VM only has a single "wan" interface.
	group := &Group{
		Name:        "test_group",
		Members:     api.SelectedMapList{"wan"},
		NoGroup:     "0",
		Sequence:    "5",
		Description: "test group",
	}

	key, err := controller.AddGroup(ctx, group)
	if err != nil {
		t.Fatalf("Failed to add group: %v", err)
	}
	t.Logf("Added group with key: %s", key)

	retrieved, err := controller.GetGroup(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get group: %v", err)
	}
	t.Logf("Retrieved group: %+v", retrieved)
	if retrieved.Name != group.Name {
		t.Fatalf("Retrieved group name does not match: got %s, want %s", retrieved.Name, group.Name)
	}
	if len(retrieved.Members) != 1 || retrieved.Members[0] != "wan" {
		t.Fatalf("Retrieved group members do not match: got %v, want [wan]", retrieved.Members)
	}
	if retrieved.Sequence != group.Sequence {
		t.Fatalf("Retrieved group sequence does not match: got %s, want %s", retrieved.Sequence, group.Sequence)
	}
	if retrieved.Description != group.Description {
		t.Fatalf("Retrieved group description does not match: got %s, want %s", retrieved.Description, group.Description)
	}

	group.Name = "test_group_upd"
	group.NoGroup = "1"
	err = controller.UpdateGroup(ctx, key, group)
	if err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}

	retrieved, err = controller.GetGroup(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated group: %v", err)
	}
	if retrieved.Name != "test_group_upd" {
		t.Fatalf("Retrieved group name does not match updated name: got %s, want %s", retrieved.Name, "test_group_upd")
	}
	if retrieved.NoGroup != "1" {
		t.Fatalf("Retrieved group nogroup does not match: got %s, want 1", retrieved.NoGroup)
	}

	err = controller.DeleteGroup(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete group: %v", err)
	}
	t.Logf("Deleted group with key: %s", key)
}
