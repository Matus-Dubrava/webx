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
		{Spath: "/api", Tpath: "/", Taddr: "127.0.0.1", Tport: 8000},
		{Spath: "/users", Tpath: "/users", Taddr: "10.0.0.1", Tport: 9000},
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
