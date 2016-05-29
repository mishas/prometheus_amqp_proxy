package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AMQPConfig struct {
	URL      string `yaml:"url,omitempty"`
	Exchange string `yaml:"exchange,omitempty"`
}

// TLSConfig configures the options for TLS connections.
type TLSConfig struct {
	// The CA cert to use for the targets.
	CAFile string `yaml:"ca_file,omitempty"`
	// The client cert file for the targets.
	CertFile string `yaml:"cert_file,omitempty"`
	// The client key file for the targets.
	KeyFile string `yaml:"key_file,omitempty"`
	// Used to verify the hostname for the targets.
	ServerName string `yaml:"server_name,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

type TargetGroup struct {
	Targets []string
}

type ScrapeConfig struct {
	// The job name to which the job label is set by default.
	JobName    string     `yaml:"job_name"`
	AMQPConfig AMQPConfig `yaml:"amqp_config,omitempty"`
	// TLSConfig to use to connect to the targets.
	TLSConfig TLSConfig `yaml:"tls_config,omitempty"`

	// List of labeled target groups for this job.
	TargetGroups []*TargetGroup `yaml:"target_groups,omitempty"`
}

type Config struct {
	ScrapeConfigs []*ScrapeConfig `yaml:"scrape_configs,omitempty"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`

	// original is the input from which the config was parsed.
	original string
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.original = s
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
