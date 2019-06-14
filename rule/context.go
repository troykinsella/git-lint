package rule

type Context struct {
	Failures []Failure
	Warnings []Failure
}

func (ctx *Context) AddFailure(f Failure, warning bool) {
	if warning {
		ctx.Warnings = append(ctx.Warnings, f)
	} else {
		ctx.Failures = append(ctx.Failures, f)
	}
}

func (ctx *Context) Passed() bool {
	return len(ctx.Failures) == 0
}
