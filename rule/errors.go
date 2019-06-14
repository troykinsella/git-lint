package rule

type ConfigError struct {
	msg string
}

func NewConfigError(msg string) *ConfigError {
	return &ConfigError{
		msg: msg,
	}
}

func (ce *ConfigError) Error() string {
	return ce.msg
}
