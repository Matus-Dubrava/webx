package core

import (
	"testing"
)

func TestConfigPassRules(t *testing.T) {
	confpath := "../testdata/configs/pass_rules_valid.toml"

	conf, err := ParseConfig(confpath)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len := len(conf.PassRules); len != 2 {
		t.Fatalf("invalid number of config rules; expected 2, got %d", len)
	}

	expected := []ProxyPassRule{
		{Spath: "/api", Tpath: "/", Thost: "127.0.0.1", Tport: 8000, HcPath: "/health"},
		{Spath: "/users", Tpath: "/users", Thost: "10.0.0.1", Tport: 9000, HcPath: "/health"},
	}

	for i := range expected {
		if expected[i] != conf.PassRules[i] {
			t.Fatalf("failed to parse proxy pass rule; expected: %+v, got: %+v", expected[i], conf.PassRules[i])
		}
	}
}

func TestConfigPassRulesInvalid1(t *testing.T) {
	confpath := "../testdata/configs/pass_rules_incomplete_rule_1.toml"
	_, err := ParseConfig(confpath)

	if err == nil {
		t.Fatalf("expected to fail to parse config file")
	}
}

func TestDuplicatePassRules(t *testing.T) {
	confpath := "../testdata/configs/duplicate_pass_rules.toml"
	_, err := ParseConfig(confpath)

	if err == nil {
		t.Fatalf("expected to fail when duplicate pass rules are detected")
	}
}

func TestValidGlobalSection(t *testing.T) {
	confpath := "../testdata/configs/valid_global_section.toml"
	conf, err := ParseConfig(confpath)
	if err != nil {
		t.Fatal(err)
	}

	expected_addr := "127.0.0.1"
	if conf.Global.Listener_host != expected_addr {
		t.Fatalf("invalid listener address; expected %s, got %s", expected_addr, conf.Global.Listener_host)
	}

	expected_port := 8080
	if conf.Global.Listener_port != expected_port {
		t.Fatalf("invalid listener port; expected %d, got %d", expected_port, conf.Global.Listener_port)
	}
}

func TestGlobalSection(t *testing.T) {
	confpath := "../testdata/configs/global_missing_fields.toml"
	_, err := ParseConfig(confpath)

	if err == nil {
		t.Fatalf("expected to fail when global fields are missing")
	}
}
