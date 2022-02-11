package api

import (
	"fmt"
	"strings"
)

// Config holds BBB api-related data.
type Config struct {
	Secret string `yaml:"secret"`
	Host   string `yaml:"host"`
}

// Sanitization check and sanitize api config instance.
func (c *Config) Sanitization() error {
	if c.Secret == "" {
		return fmt.Errorf("`secret` field is required")
	}

	if c.Host == "" {
		return fmt.Errorf("`host` field is required")
	}

	if !strings.HasSuffix(c.Host, "/") {
		c.Host += "/"
	}

	return nil
}
