package rule

type Rule interface {
	ID() string
	Name() string
	SetConfig(config interface{}) error
	GetConfig() Config
	Run(ctx *Context) error
}
