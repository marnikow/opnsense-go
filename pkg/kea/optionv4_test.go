package kea

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestOptionV4(t *testing.T) {
	controller := newDdnsTestController()
	ctx := context.Background()

	option := &OptionV4{
		Code:        "43",
		Encoding:    "hex",
		Data:        "0104c0000201",
		Force:       "0",
		Description: "test option 43",
	}
	key, err := controller.AddOptionV4(ctx, option)
	if err != nil {
		t.Fatalf("Failed to add OptionV4: %v", err)
	}
	t.Logf("Added OptionV4 with key: %s", key)
	defer func() {
		if err := controller.DeleteOptionV4(ctx, key); err != nil {
			t.Errorf("Failed to delete OptionV4: %v", err)
		}
	}()

	got, err := controller.GetOptionV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get OptionV4: %v", err)
	}
	if got.Code != "43" || got.Encoding != "hex" || got.Data != option.Data || got.Description != option.Description {
		t.Fatalf("Retrieved OptionV4 does not match: %+v", got)
	}

	// Reference the option from a subnet.
	subnet := &SubnetV4{
		Subnet:                "192.168.202.0/24",
		MatchClientId:         "1",
		OptionDataAutoCollect: "1",
		OptionData: OptionDataV4{
			DomainNameServers: api.SelectedMapList{},
			Routers:           api.SelectedMapList{"192.168.202.1"},
			NtpServers:        api.SelectedMapList{},
			TimeServers:       api.SelectedMapList{},
			DomainSearch:      api.SelectedMapList{},
		},
		Option:      api.SelectedMapList{key},
		Description: "Test subnet for OptionV4",
	}
	subnetKey, err := controller.AddSubnetV4(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV4 with option: %v", err)
	}
	defer func() {
		if err := controller.DeleteSubnetV4(ctx, subnetKey); err != nil {
			t.Errorf("Failed to delete SubnetV4: %v", err)
		}
	}()
	gotSubnet, err := controller.GetSubnetV4(ctx, subnetKey)
	if err != nil {
		t.Fatalf("Failed to get SubnetV4: %v", err)
	}
	if len(gotSubnet.Option) != 1 || gotSubnet.Option[0] != key {
		t.Fatalf("SubnetV4 option reference does not match: got %v, want [%s]", gotSubnet.Option, key)
	}

	// Allocator and ping check round-trip on the same subnet.
	subnet.Allocator = "random"
	subnet.PingCheck = "1"
	if err := controller.UpdateSubnetV4(ctx, subnetKey, subnet); err != nil {
		t.Fatalf("Failed to update SubnetV4: %v", err)
	}
	gotSubnet, err = controller.GetSubnetV4(ctx, subnetKey)
	if err != nil {
		t.Fatalf("Failed to get updated SubnetV4: %v", err)
	}
	if gotSubnet.Allocator != "random" || gotSubnet.PingCheck != "1" {
		t.Fatalf("SubnetV4 allocator/ping_check not updated: %+v", gotSubnet)
	}

	// Lifetime, client class match and the newer option data fields.
	subnet.ValidLifetime = "7200"
	subnet.MatchCode = "60"
	subnet.MatchEncoding = "string"
	subnet.MatchData = "ubnt"
	subnet.OptionData.ClasslessStaticRoute = "10.10.0.0/16 - 192.168.202.1"
	subnet.OptionData.V6OnlyPreferred = "300"
	if err := controller.UpdateSubnetV4(ctx, subnetKey, subnet); err != nil {
		t.Fatalf("Failed to update SubnetV4 (lifetime/match): %v", err)
	}
	gotSubnet, err = controller.GetSubnetV4(ctx, subnetKey)
	if err != nil {
		t.Fatalf("Failed to get updated SubnetV4: %v", err)
	}
	if gotSubnet.ValidLifetime != "7200" || gotSubnet.MatchCode != "60" || gotSubnet.MatchData != "ubnt" ||
		gotSubnet.OptionData.ClasslessStaticRoute == "" || gotSubnet.OptionData.V6OnlyPreferred != "300" || gotSubnet.SubnetId == "" {
		t.Fatalf("SubnetV4 lifetime/match/option data not updated: %+v", gotSubnet)
	}

	// Match condition: only send when the client sends vendor class 60.
	option.MatchCode = "60"
	option.MatchEncoding = "string"
	option.MatchData = "ubnt"
	option.Force = "1"
	if err := controller.UpdateOptionV4(ctx, key, option); err != nil {
		t.Fatalf("Failed to update OptionV4: %v", err)
	}
	got, err = controller.GetOptionV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated OptionV4: %v", err)
	}
	if got.MatchCode != "60" || got.MatchEncoding != "string" || got.MatchData != "ubnt" || got.Force != "1" {
		t.Fatalf("Updated OptionV4 does not match: %+v", got)
	}
}
