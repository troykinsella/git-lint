package branchnaming

import (
	"regexp"
)

type Config struct {
	Warn bool
	Allow bool
	Patterns []string

	res []*regexp.Regexp
}

func (c *Config) IsWarn() bool {
	return c.Warn
}
