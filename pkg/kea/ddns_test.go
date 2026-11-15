package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func newDdnsTestController() Controller {
	return Controller{Api: api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})}
}

func TestSettingsDdns(t *testing.T) {
	controller := newDdnsTestController()
	ctx := context.Background()

	respGet, err := controller.SettingsDdnsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsDdnsGet failed: %v", err)
	}
	t.Logf("SettingsDdnsGet (initial): %+v", respGet)
	if respGet.Ddns.General.ServerIp == "" || respGet.Ddns.General.ServerPort == "" {
		t.Fatal("Expected general.server_ip and general.server_port to be set")
	}

	orig := respGet.Ddns

	updated := orig
	updated.General.ServerPort = "53002"
	respSet, err := controller.SettingsDdnsUpdate(ctx, &updated)
	if err != nil {
		t.Fatalf("SettingsDdnsUpdate failed: %v", err)
	}
	if respSet.Result != "saved" {
		t.Fatalf("Expected result 'saved', got %q", respSet.Result)
	}

	respGet, err = controller.SettingsDdnsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsDdnsGet (after update) failed: %v", err)
	}
	if respGet.Ddns.General.ServerPort != "53002" {
		t.Fatalf("server_port not updated; got %q, want %q", respGet.Ddns.General.ServerPort, "53002")
	}

	if _, err := controller.SettingsDdnsUpdate(ctx, &orig); err != nil {
		t.Fatalf("SettingsDdnsUpdate (restore) failed: %v", err)
	}
	respGet, err = controller.SettingsDdnsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsDdnsGet (after restore) failed: %v", err)
	}
	if respGet.Ddns.General.ServerPort != orig.General.ServerPort {
		t.Fatalf("server_port not restored; got %q, want %q", respGet.Ddns.General.ServerPort, orig.General.ServerPort)
	}
}

func TestSubnetV4Ddns(t *testing.T) {
	controller := newDdnsTestController()
	ctx := context.Background()

	subnet := &SubnetV4{
		Subnet:                     "192.168.201.0/24",
		Pools:                      "192.168.201.100 - 192.168.201.200",
		MatchClientId:              "1",
		OptionDataAutoCollect:      "0",
		DdnsForwardZone:            "test.example.com",
		DdnsReverseZone:            "201.168.192.in-addr.arpa",
		DdnsQualifyingSuffix:       "test.example.com",
		DdnsDnsServer:              "127.0.0.1",
		DdnsDnsPort:                "53530",
		DdnsDomainKeyName:          "test-key",
		DdnsDomainKeySecret:        "c2VjcmV0LXNlY3JldC1zZWNyZXQ=",
		DdnsDomainKeyAlgorithm:     api.SelectedMap("hmac-sha512"),
		DdnsOverrideNoUpdate:       "1",
		DdnsOverrideClientUpdate:   "1",
		DdnsUpdateOnRenew:          "0",
		DdnsConflictResolutionMode: api.SelectedMap("no-check-with-dhcid"),
		Description:                "Test Kea DHCPv4 Subnet DDNS",
	}

	key, err := controller.AddSubnetV4(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV4: %v", err)
	}
	t.Cleanup(func() {
		if err := controller.DeleteSubnetV4(ctx, key); err != nil {
			t.Errorf("Failed to delete SubnetV4: %v", err)
		}
	})

	got, err := controller.GetSubnetV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get SubnetV4: %v", err)
	}

	checks := map[string][2]string{
		"ddns_forward_zone":             {got.DdnsForwardZone, subnet.DdnsForwardZone},
		"ddns_reverse_zone":             {got.DdnsReverseZone, subnet.DdnsReverseZone},
		"ddns_qualifying_suffix":        {got.DdnsQualifyingSuffix, subnet.DdnsQualifyingSuffix},
		"ddns_dns_server":               {got.DdnsDnsServer, subnet.DdnsDnsServer},
		"ddns_dns_port":                 {got.DdnsDnsPort, subnet.DdnsDnsPort},
		"ddns_domain_key_name":          {got.DdnsDomainKeyName, subnet.DdnsDomainKeyName},
		"ddns_domain_key_secret":        {got.DdnsDomainKeySecret, subnet.DdnsDomainKeySecret},
		"ddns_domain_key_algorithm":     {string(got.DdnsDomainKeyAlgorithm), string(subnet.DdnsDomainKeyAlgorithm)},
		"ddns_override_no_update":       {got.DdnsOverrideNoUpdate, subnet.DdnsOverrideNoUpdate},
		"ddns_override_client_update":   {got.DdnsOverrideClientUpdate, subnet.DdnsOverrideClientUpdate},
		"ddns_update_on_renew":          {got.DdnsUpdateOnRenew, subnet.DdnsUpdateOnRenew},
		"ddns_conflict_resolution_mode": {string(got.DdnsConflictResolutionMode), string(subnet.DdnsConflictResolutionMode)},
	}
	for field, v := range checks {
		if v[0] != v[1] {
			t.Errorf("%s mismatch: got %q, want %q", field, v[0], v[1])
		}
	}

	// Clearing the DDNS settings must round-trip as empty values.
	subnet.DdnsForwardZone = ""
	subnet.DdnsReverseZone = ""
	subnet.DdnsQualifyingSuffix = ""
	subnet.DdnsDnsServer = ""
	subnet.DdnsDnsPort = ""
	subnet.DdnsDomainKeyName = ""
	subnet.DdnsDomainKeySecret = ""
	subnet.DdnsDomainKeyAlgorithm = ""
	subnet.DdnsConflictResolutionMode = ""
	if err := controller.UpdateSubnetV4(ctx, key, subnet); err != nil {
		t.Fatalf("Failed to update SubnetV4: %v", err)
	}
	got, err = controller.GetSubnetV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get SubnetV4 after update: %v", err)
	}
	if got.DdnsForwardZone != "" || got.DdnsDomainKeySecret != "" || got.DdnsConflictResolutionMode != "" {
		t.Errorf("DDNS fields not cleared: %+v", got)
	}
}
