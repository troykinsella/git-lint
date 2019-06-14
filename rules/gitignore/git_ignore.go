package gitignore

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/rule"
	"io/ioutil"
	"strings"
)

const (
	ID   = "GL005"
	Name = "git_ignore"
)

type GitIgnore struct {
	config *Config
}

func (gi *GitIgnore) ID() string {
	return ID
}

func (gi *GitIgnore) Name() string {
	return Name
}

func (gi *GitIgnore) SetConfig(rawConfig interface{}) error {
	var config Config
	err := mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	gi.config = &config

	return nil
}

func (gi *GitIgnore) GetConfig() rule.Config {
	return gi.config
}

func (gi *GitIgnore) Run(ctx *rule.Context) error {

	fBytes, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return err
	}

	ignoredEntries := strings.Split(string(fBytes), "\n")

	var found bool
	for _, expectedEntry := range gi.config.Entries {
		found = false
		for _, ignoredEntry := range ignoredEntries {
			if ignoredEntry == expectedEntry {
				found = true
				break
			}
		}

		if !found {
			ctx.AddFailure(rule.NewBasicFailure(
				gi,
				fmt.Sprintf("expected git to ignore: '%s'", expectedEntry),
			), gi.config.IsWarn())
		}
	}

	return nil
}
