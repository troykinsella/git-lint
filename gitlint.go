package git_lint

import (
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
)

var (
	AppVersion = "0.0.0-dev.0"
)

type GitLint struct {
	config *Config
	repo   *git.Repository
}

func New(config *Config) (*GitLint, error) {
	if config == nil {
		panic("config required")
	}

	repo, err := git.NewRepository("", "origin")
	if err != nil {
		return nil, err
	}

	return &GitLint{
		config: config,
		repo:   repo,
	}, nil
}

func (gl *GitLint) loadRules() ([]rule.Rule, error) {
	rules := make([]rule.Rule, 0)

	for name, ruleConfig := range gl.config.Rules {
		rule := NewRule(name, gl.repo)
		err := rule.SetConfig(ruleConfig)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func (gl *GitLint) Run() (bool, error) {

	ctx := &rule.Context{}

	rules, err := gl.loadRules()
	if err != nil {
		return false, err
	}

	for _, rule := range rules {
		err = rule.Run(ctx)
		if err != nil {
			return false, err
		}
	}

	reporter := NewReporter(ctx)
	reporter.Print()

	return ctx.Passed(), nil
}
