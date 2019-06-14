package git_lint

type Config struct {
	Version string                 `yaml:"version"`
	Rules   map[string]interface{} `yaml:"rules"`
}
