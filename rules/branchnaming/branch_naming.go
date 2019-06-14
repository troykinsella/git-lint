package branchnaming

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"github.com/troykinsella/git-lint/util"
	"strings"
)

const (
	ID = "GL002"
	Name = "branch_naming"
)

type BranchNaming struct {
	Repo *git.Repository
	config *Config
}

func (bn *BranchNaming) ID() string {
	return ID
}

func (bn *BranchNaming) Name() string {
	return Name
}

func (bn *BranchNaming) SetConfig(rawConfig interface{}) error {
	var config Config
	var err error

	err = mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	if len(config.Patterns) == 0 {
		return rule.NewConfigError(Name + ": must supply at least one pattern")
	}

	config.res, err = util.CompileRegexes(config.Patterns)
	if err != nil {
		return rule.NewConfigError(Name + ": invalid pattern: '" + err.Error() + "'")
	}

	bn.config = &config

	return nil
}

func (bc *BranchNaming) GetConfig() rule.Config {
	return bc.config
}

func (bn *BranchNaming) Run(ctx *rule.Context) error {
	branches, err := bn.Repo.Branches(&git.BranchOpts{
		ShortNames: true,
	})
	if err != nil {
		return err
	}

	var matched bool

	for _, branch := range branches {
		matched = false

		for _, re := range bn.config.res {
			if re.Match([]byte(branch)) {
				matched = true

				if !bn.config.Allow {
					ctx.AddFailure(rule.NewRegexpFailure(
						bn,
						"branch name matched a disallowed pattern",
						re.String(),
						branch,
						false,
					), bn.config.IsWarn())
					break
				}
			}
		}

		if bn.config.Allow && !matched {
			ctx.AddFailure(rule.NewRegexpFailure(
				bn,
				"branch name did not match allowed patterns",
				fmt.Sprintf("/%s/", strings.Join(bn.config.Patterns, "/, /")),
				branch,
				true,
			), bn.config.IsWarn())
		}
	}

	return nil
}
