package tagname

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"github.com/troykinsella/git-lint/util"
	"strings"
)

const (
	ID   = "GL007"
	Name = "tag_name"
)

type TagName struct {
	Repo   *git.Repository
	config *Config
}

func (tn *TagName) ID() string {
	return ID
}

func (tn *TagName) Name() string {
	return Name
}

func (tn *TagName) SetConfig(rawConfig interface{}) error {
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

	tn.config = &config

	return nil
}

func (tn *TagName) GetConfig() rule.Config {
	return tn.config
}

func (tn *TagName) Run(ctx *rule.Context) error {
	tags, err := tn.Repo.Tags()
	if err != nil {
		return err
	}

	var matched bool

	for _, tag := range tags {
		for _, re := range tn.config.res {
			if re.Match([]byte(tag)) {
				matched = true

				if !tn.config.Allow {
					ctx.AddFailure(rule.NewRegexpFailure(
						tn,
						"tag name matched a disallowed pattern",
						re.String(),
						tag,
						false,
					), tn.config.IsWarn())
					break
				}
			}
		}

		if tn.config.Allow && !matched {
			ctx.AddFailure(rule.NewRegexpFailure(
				tn,
				"tag name did not match allowed patterns",
				fmt.Sprintf("/%s/", strings.Join(tn.config.Patterns, "/, /")),
				tag,
				true,
			), tn.config.IsWarn())
		}
	}

	return nil
}
