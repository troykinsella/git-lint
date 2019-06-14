package gitkeep

type Config struct {
	Warn        bool
	Directories []string
}

func (c *Config) IsWarn() bool {
	return c.Warn
}
