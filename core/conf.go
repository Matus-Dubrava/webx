package core

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type ProxyPassRule struct {
	Spath string `toml:"source_path"`
	Tpath string `toml:"target_path"`
	Thost string `toml:"target_host"`
	Tport int    `toml:"target_port"`
}

func (rule *ProxyPassRule) ToString() string {
	return fmt.Sprintf("source_path: %s, target_path: %s, target_host: %s, target_port: %d", rule.Spath, rule.Tpath, rule.Thost, rule.Tport)
}

type ConfigGlobal struct {
	Listener_host string `toml:"listener_host"`
	Listener_port int    `toml:"listener_port"`
}

type Config struct {
	Global    ConfigGlobal    `toml:"global"`
	PassRules []ProxyPassRule `toml:"proxy_pass"`
}

func ValidateGlobalSection(conf *Config) error {
	if conf.Global.Listener_host == "" {
		return fmt.Errorf("invalid config: missing 'listener_host' field")
	}

	if conf.Global.Listener_port == 0 {
		return fmt.Errorf("invalid config: mssing 'listener_port' field or port is 0")
	}

	return nil
}

func ValidatePassRule(rule *ProxyPassRule) error {
	if rule.Spath == "" {
		return fmt.Errorf("invalid proxy_pass: missing 'source_path' field; rule: %s", rule.ToString())
	}
	if rule.Tpath == "" {
		return fmt.Errorf("invalid proxy_pass: missing 'target_path' field; rule: %s", rule.ToString())
	}
	if rule.Thost == "" {
		return fmt.Errorf("invalid proxy_pass: missing 'target_host' field; rule: %s", rule.ToString())
	}
	if rule.Tport == 0 {
		return fmt.Errorf("invalid proxy_pass: missing 'target_port' field or field is 0; rule: %s", rule.ToString())
	}

	return nil
}

func ValidatePassRules(conf *Config) error {
	paths := make(map[string]struct{})

	for _, rule := range conf.PassRules {
		// Ensure that there are no multiple rules with the same source_path
		if _, exists := paths[rule.Spath]; exists {
			return fmt.Errorf("duplicate source path in rule %s", rule.ToString())
		}

		paths[rule.Spath] = struct{}{}

		if err := ValidatePassRule(&rule); err != nil {
			return err
		}
	}

	return nil
}

func ParseConfig(filepath string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(filepath, &conf); err != nil {
		return &Config{}, err
	}

	if err := ValidatePassRules(&conf); err != nil {
		return &Config{}, err
	}

	if err := ValidateGlobalSection(&conf); err != nil {
		return &Config{}, err
	}

	return &conf, nil
}

func (conf *Config) Print() {
	if conf == nil {
		log.Println("error: can' print config; config is nil")
		return
	}

	fmt.Println("[config]")
	fmt.Printf("[global.listener_host] = %s\n", conf.Global.Listener_host)
	fmt.Printf("[global.listener_port] = %d\n", conf.Global.Listener_port)
	fmt.Println("[proxy pass rules]")
	for i, rule := range conf.PassRules {
		fmt.Printf("rule %d: %s\n", i+1, rule.ToString())
	}
}
