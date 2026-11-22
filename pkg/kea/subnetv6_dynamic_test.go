package kea

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// TestSubnetV6AutoCollect round-trips option_data_autocollect and
// dynamic_prefix. A dynamic prefix subnet needs an interface with an
// identity association prefix, which the CI VM does not have, so only the
// static case (dynamic_prefix = 0) is created here.
func TestSubnetV6AutoCollect(t *testing.T) {
	controller := newDdnsTestController()
	ctx := context.Background()

	subnet := &SubnetV6{
		Subnet:                "fd00:201::/64",
		Interface:             api.SelectedMap("wan"),
		OptionData:            OptionDataV6{DomainNameServers: api.SelectedMapList{}, DomainSearch: api.SelectedMapList{}},
		DynamicPrefix:         "0",
		OptionDataAutoCollect: "1",
		Description:           "Test Kea DHCPv6 autocollect",
	}
	key, err := controller.AddSubnetV6(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV6: %v", err)
	}
	defer func() {
		if err := controller.DeleteSubnetV6(ctx, key); err != nil {
			t.Errorf("Failed to delete SubnetV6: %v", err)
		}
	}()

	got, err := controller.GetSubnetV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get SubnetV6: %v", err)
	}
	if got.OptionDataAutoCollect != "1" || got.DynamicPrefix != "0" {
		t.Fatalf("Retrieved SubnetV6 does not match: %+v", got)
	}
}
