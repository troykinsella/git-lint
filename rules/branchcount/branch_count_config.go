package branchcount

type Config struct {
	Warn bool
	Max  int
}

func (c *Config) IsWarn() bool {
	return c.Warn
}
