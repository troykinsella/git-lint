package gitignore

type Config struct {
	Warn    bool
	Entries []string
}

func (c *Config) IsWarn() bool {
	return c.Warn
}
