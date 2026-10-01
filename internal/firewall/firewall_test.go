package firewall

import (
	"context"
	"testing"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func rule(port, proto, typ string) elestio.ServiceFirewallRule {
	return elestio.ServiceFirewallRule{Port: port, Protocol: proto, Type: typ, Targets: GetDefaultTargets()}
}

func TestMergeWithToolPorts(t *testing.T) {
	user := []elestio.ServiceFirewallRule{rule("443", "tcp", "INPUT")}
	api := []elestio.ServiceFirewallRule{
		rule("443", "tcp", "INPUT"),
		rule("18345", "tcp", "INPUT"), // VS Code tool port managed by the API
		rule("9999", "tcp", "INPUT"),  // unknown port: must not be resurrected
	}

	t.Run("keeps API tool ports and drops unknown API rules", func(t *testing.T) {
		got := MergeWithToolPorts(user, api, false)
		ports := map[string]int{}
		for _, r := range got {
			ports[r.Port]++
		}
		if ports["443"] != 1 || ports["18345"] != 1 || ports["9999"] != 0 || len(got) != 2 {
			t.Fatalf("unexpected merge result: %+v", got)
		}
	})

	t.Run("remove_tool_ports drops tool ports", func(t *testing.T) {
		got := MergeWithToolPorts(user, api, true)
		if len(got) != 1 || got[0].Port != "443" {
			t.Fatalf("unexpected merge result: %+v", got)
		}
	})

	t.Run("does not mutate its inputs", func(t *testing.T) {
		u := []elestio.ServiceFirewallRule{rule("80", "tcp", "INPUT")}
		_ = MergeWithToolPorts(u, api, false)
		if len(u) != 1 {
			t.Fatal("user rules slice was modified")
		}
	})
}

func TestIsToolPort(t *testing.T) {
	if !IsToolPort("18345", "tcp", "input") {
		t.Error("18345/tcp/input is a tool port")
	}
	for _, c := range [][3]string{{"18345", "udp", "input"}, {"22", "tcp", "input"}, {"18345", "tcp", "output"}} {
		if IsToolPort(c[0], c[1], c[2]) {
			t.Errorf("%v must not be a tool port", c)
		}
	}
}

func TestRequiredSystemPortsPresent(t *testing.T) {
	got := GetRequiredSystemPorts()
	seen := map[string]bool{}
	for _, r := range got {
		seen[r.Port+"/"+r.Protocol] = true
	}
	// SSH must always stay reachable, and Nebula is the management overlay.
	if !seen["22/tcp"] || !seen["4242/udp"] {
		t.Fatalf("required system ports missing: %+v", got)
	}
}

func TestConvertTerraformToElestio_NullSetIsEmpty(t *testing.T) {
	var diags diag.Diagnostics
	set := ConvertElestioToTerraform(nil, &diags)
	rules := ConvertTerraformToElestio(context.Background(), set, &diags)
	if diags.HasError() || len(rules) != 0 {
		t.Fatalf("rules=%v diags=%v", rules, diags)
	}
}

func TestDefaultTargetsAreNotMutableShared(t *testing.T) {
	a := GetDefaultTargets()
	a[0] = "10.0.0.0/8"
	if GetDefaultTargets()[0] == "10.0.0.0/8" {
		t.Fatal("GetDefaultTargets returns shared mutable state")
	}
}
