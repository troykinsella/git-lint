package branchlastcommit

import (
	"github.com/troykinsella/sloppy-duration"
)

type Config struct {
	Warn         bool
	Max_Duration string

	maxDuration *sloppy_duration.SloppyDuration
}

func (c *Config) IsWarn() bool {
	return c.Warn
}
