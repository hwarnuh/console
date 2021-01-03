package git

import (
	"flag"
	"fmt"
	"time"
)

// Config for Git Service
type Config struct {
	Enabled bool `yaml:"enabled"`

	// AllowedFileExtensions specifies file extensions that shall be picked up. If at least one is specified all other
	// file extensions will be ignored.
	AllowedFileExtensions []string `yaml:"-"`

	// Max file size which will be considered. Files exceeding this size will be ignored and logged.
	MaxFileSize int64 `yaml:"-"`

	// RefreshInterval specifies how often the repository shall be pulled to check for new changes.
	RefreshInterval time.Duration `yaml:"refreshInterval"`

	// Repository that contains markdown files that document a Kafka topic.
	Repository RepositoryConfig `yaml:"repository"`

	// Authentication Configs
	BasicAuth BasicAuthConfig `yaml:"basicAuth"`
	SSH       SSHConfig       `yaml:"ssh"`
}

// RegisterFlags for all (sub)configs
func (c *Config) RegisterFlags(f *flag.FlagSet) {
	c.BasicAuth.RegisterFlags(f)
	c.SSH.RegisterFlags(f)
}

// Validate all root and child config structs
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.RefreshInterval == 0 {
		return fmt.Errorf("git config is enabled but refresh interval is set to 0 (disabled)")
	}

	return c.Repository.Validate()
}

// SetDefaults for all root and child config structs
func (c *Config) SetDefaults() {
	c.RefreshInterval = time.Minute
	c.MaxFileSize = 500 * 1000 // 500KB
}
