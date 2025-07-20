package core

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type ProxyPassRule struct {
	Spath string `toml:"source_path"`
	Tpath string `toml:"target_path"`
	Taddr string `toml:"target_addr"`
	Tport int    `toml:"target_port"`
}

func (rule *ProxyPassRule) ToString() string {
	return fmt.Sprintf("source_path: %s, target_path: %s, target_addr: %s, target_port: %d", rule.Spath, rule.Tpath, rule.Taddr, rule.Tport)
}

type ConfigGlobal struct {
	Listener_addr string `toml:"listener_address"`
	Listener_port int    `toml:"listener_port"`
}

type Config struct {
	Global    ConfigGlobal    `toml:"global"`
	PassRules []ProxyPassRule `toml:"proxy_pass"`
}

func ValidatePassRule(rule *ProxyPassRule) error {
	if rule.Spath == "" {
		return fmt.Errorf("invalid proxy_pass: missing 'source_path' field; rule: %s", rule.ToString())
	}
	if rule.Tpath == "" {
		return fmt.Errorf("invalid proxy_pass: missing 'target_path' field; rule: %s", rule.ToString())
	}
	if rule.Taddr == "" {
		return fmt.Errorf("invalid proxy_pass: missing 'target_addr' field; rule: %s", rule.ToString())
	}
	if rule.Tport == 0 {
		return fmt.Errorf("invalid proxy_pass: missing 'target_port' field or field is 0; rule: %s", rule.ToString())
	}

	return nil
}

func ParseConfig(filepath string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(filepath, &conf); err != nil {
		return &Config{}, err
	}

	for _, rule := range conf.PassRules {
		if err := ValidatePassRule(&rule); err != nil {
			return &Config{}, err
		}
	}

	return &conf, nil
}

func (conf *Config) Print() {
	if conf == nil {
		log.Println("error: can' print config; config is nil")
		return
	}

	fmt.Println("[config]")
	fmt.Printf("[global.listener_addr] = %s\n", conf.Global.Listener_addr)
	fmt.Printf("[global.listener_port] = %d\n", conf.Global.Listener_port)
	fmt.Println("[proxy pass rules]")
	for i, rule := range conf.PassRules {
		fmt.Printf("rule %d: %s\n", i+1, rule.ToString())
	}
}
