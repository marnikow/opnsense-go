package kea

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestSettingsGeneralFields round-trips the general fields that are not
// covered by TestSettingsDHCPv4/TestSettingsDHCPv6, then restores them.
func TestSettingsGeneralFields(t *testing.T) {
	controller := newDdnsTestController()
	ctx := context.Background()

	v4, err := controller.SettingsDHCPv4Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Get failed: %v", err)
	}
	orig4 := v4.DHCPv4
	upd4 := orig4
	upd4.General.ServiceSocketsMaxRetries = "5"
	upd4.General.ServiceSocketsRetryWaitTime = "2000"
	upd4.General.DeclineProbationPeriod = "600"
	upd4.General.Compatibility = api.SelectedMapList{"ignore-rai-link-selection"}
	if _, err := controller.SettingsDHCPv4Update(ctx, &upd4); err != nil {
		t.Fatalf("SettingsDHCPv4Update failed: %v", err)
	}
	v4, err = controller.SettingsDHCPv4Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Get (after update) failed: %v", err)
	}
	g4 := v4.DHCPv4.General
	if g4.ServiceSocketsMaxRetries != "5" || g4.ServiceSocketsRetryWaitTime != "2000" || g4.DeclineProbationPeriod != "600" ||
		len(g4.Compatibility) != 1 || g4.Compatibility[0] != "ignore-rai-link-selection" {
		t.Fatalf("DHCPv4 general fields not updated: %+v", g4)
	}
	if _, err := controller.SettingsDHCPv4Update(ctx, &orig4); err != nil {
		t.Fatalf("SettingsDHCPv4Update (restore) failed: %v", err)
	}

	v6, err := controller.SettingsDHCPv6Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Get failed: %v", err)
	}
	orig6 := v6.DHCPv6
	upd6 := orig6
	upd6.General.DeclineProbationPeriod = "600"
	upd6.General.MacSources = api.SelectedMapList{"duid", "ipv6-link-local"}
	if _, err := controller.SettingsDHCPv6Update(ctx, &upd6); err != nil {
		t.Fatalf("SettingsDHCPv6Update failed: %v", err)
	}
	v6, err = controller.SettingsDHCPv6Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Get (after update) failed: %v", err)
	}
	g6 := v6.DHCPv6.General
	if g6.DeclineProbationPeriod != "600" || len(g6.MacSources) != 2 {
		t.Fatalf("DHCPv6 general fields not updated: %+v", g6)
	}
	if _, err := controller.SettingsDHCPv6Update(ctx, &orig6); err != nil {
		t.Fatalf("SettingsDHCPv6Update (restore) failed: %v", err)
	}
}
