package core

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestTunable(t *testing.T) {
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

	// net.inet.ip.random_id = 1 is the OPNsense default, so applying it is harmless.
	tunable := &Tunable{
		Tunable:     "net.inet.ip.random_id",
		Value:       "1",
		Description: "test tunable",
	}
	key, err := controller.AddTunable(ctx, tunable)
	if err != nil {
		t.Fatalf("Failed to add tunable: %v", err)
	}
	t.Logf("Added tunable with key: %s", key)
	defer func() {
		if err := controller.DeleteTunable(ctx, key); err != nil {
			t.Errorf("Failed to delete tunable: %v", err)
		}
	}()

	got, err := controller.GetTunable(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get tunable: %v", err)
	}
	if got.Tunable != tunable.Tunable || got.Value != "1" || got.Description != tunable.Description {
		t.Fatalf("Retrieved tunable does not match: %+v", got)
	}

	tunable.Description = "test tunable updated"
	if err := controller.UpdateTunable(ctx, key, tunable); err != nil {
		t.Fatalf("Failed to update tunable: %v", err)
	}
	got, err = controller.GetTunable(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated tunable: %v", err)
	}
	if got.Description != "test tunable updated" {
		t.Fatalf("Updated tunable does not match: %+v", got)
	}
}
