package branchsingleton

import (
	"regexp"
)

type Config struct {
	Warn       bool
	Singletons []string

	res []*regexp.Regexp
}

func (c *Config) IsWarn() bool {
	return c.Warn
}
