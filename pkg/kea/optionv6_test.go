package kea

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestOptionV6(t *testing.T) {
	controller := newDdnsTestController()
	ctx := context.Background()

	// Option 23: DNS recursive name servers.
	option := &OptionV6{
		Code:        "23",
		Encoding:    "ipv6-address",
		Data:        "2001:db8::53",
		Force:       "0",
		Description: "test option 23",
	}
	key, err := controller.AddOptionV6(ctx, option)
	if err != nil {
		t.Fatalf("Failed to add OptionV6: %v", err)
	}
	t.Logf("Added OptionV6 with key: %s", key)
	defer func() {
		if err := controller.DeleteOptionV6(ctx, key); err != nil {
			t.Errorf("Failed to delete OptionV6: %v", err)
		}
	}()

	got, err := controller.GetOptionV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get OptionV6: %v", err)
	}
	if got.Code != "23" || got.Encoding != "ipv6-address" || got.Data != option.Data || got.Description != option.Description {
		t.Fatalf("Retrieved OptionV6 does not match: %+v", got)
	}

	// Reference the option from a static subnet.
	subnet := &SubnetV6{
		Subnet:        "fd00:202::/64",
		Interface:     api.SelectedMap("wan"),
		OptionData:    OptionDataV6{DomainNameServers: api.SelectedMapList{}, DomainSearch: api.SelectedMapList{}},
		DynamicPrefix: "0",
		Option:        api.SelectedMapList{key},
		Description:   "Test subnet for OptionV6",
	}
	subnetKey, err := controller.AddSubnetV6(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV6 with option: %v", err)
	}
	defer func() {
		if err := controller.DeleteSubnetV6(ctx, subnetKey); err != nil {
			t.Errorf("Failed to delete SubnetV6: %v", err)
		}
	}()
	gotSubnet, err := controller.GetSubnetV6(ctx, subnetKey)
	if err != nil {
		t.Fatalf("Failed to get SubnetV6: %v", err)
	}
	if len(gotSubnet.Option) != 1 || gotSubnet.Option[0] != key {
		t.Fatalf("SubnetV6 option reference does not match: got %v, want [%s]", gotSubnet.Option, key)
	}

	option.Force = "1"
	option.Description = "test option 23 updated"
	if err := controller.UpdateOptionV6(ctx, key, option); err != nil {
		t.Fatalf("Failed to update OptionV6: %v", err)
	}
	got, err = controller.GetOptionV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated OptionV6: %v", err)
	}
	if got.Force != "1" || got.Description != "test option 23 updated" {
		t.Fatalf("Updated OptionV6 does not match: %+v", got)
	}
}
