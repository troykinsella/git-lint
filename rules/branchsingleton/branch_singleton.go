package branchsingleton

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"github.com/troykinsella/git-lint/util"
	"regexp"
	"strings"
)

const (
	ID   = "GL004"
	Name = "branch_singleton"
)

type BranchSingleton struct {
	Repo   *git.Repository
	config *Config
}

func (bs *BranchSingleton) ID() string {
	return ID
}

func (bs *BranchSingleton) Name() string {
	return Name
}

func (bs *BranchSingleton) SetConfig(rawConfig interface{}) error {
	var config Config
	var err error

	err = mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	if len(config.Singletons) == 0 {
		return rule.NewConfigError(Name + ": must supply at least one singleton pattern")
	}

	config.res, err = util.CompileRegexes(config.Singletons)
	if err != nil {
		return rule.NewConfigError(Name + ": invalid pattern: '" + err.Error() + "'")
	}

	bs.config = &config

	return nil
}

func (bs *BranchSingleton) GetConfig() rule.Config {
	return bs.config
}

func (bs *BranchSingleton) Run(ctx *rule.Context) error {
	branches, err := bs.Repo.Branches(&git.BranchOpts{
		ShortNames: true,
	})
	if err != nil {
		return err
	}

	for _, singleton := range bs.config.res {
		matches := extractMatches(branches, singleton)
		if len(matches) > 1 {
			ctx.AddFailure(
				rule.NewActualExpectedFailure(
					bs,
					"expected branch singleton",
					strings.Join(matches, ", "),
					fmt.Sprintf("one of '%s'", singleton),
				),
				bs.config.IsWarn())
		}
	}

	return nil
}

func extractMatches(list []string, re *regexp.Regexp) []string {
	result := make([]string, 0)

	for _, e := range list {
		if re.Match([]byte(e)) {
			result = append(result, e)
		}
	}

	return result
}
