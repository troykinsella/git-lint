package git_lint

import (
	"fmt"
	"github.com/troykinsella/git-lint/rule"
)

type Reporter struct {
	ctx *rule.Context
}

func NewReporter(ctx *rule.Context) *Reporter {
	return &Reporter{
		ctx: ctx,
	}
}

func (r *Reporter) Print() {

	warnings := len(r.ctx.Warnings)
	failures := len(r.ctx.Failures)

	if warnings > 0 {
		fmt.Printf("\nWarnings (%d):\n", warnings)

		for i, w := range r.ctx.Warnings {
			fmt.Printf("%d) Warning: %s\n", i+1, w)
		}
	}

	if failures > 0 {
		fmt.Printf("\nFailures (%d):\n", failures)

		for i, f := range r.ctx.Failures {
			fmt.Printf("%d) Failure: %s\n", i+1, f)
		}
	}

	var status string
	if r.ctx.Passed() {
		status = "OK"
	} else {
		status = "FAILED"
	}

	fmt.Printf("\n[%s] %d warnings, %d failures\n",
		status,
		warnings,
		failures,
	)
}
