package gitignores

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"io/ioutil"
	"strings"
)

const (
	ID = "GL005"
	Name = "git_ignores"
)

type GitIgnores struct {
	Repo *git.Repository
	config *Config
}

func (bc *GitIgnores) ID() string {
	return ID
}

func (bc *GitIgnores) Name() string {
	return Name
}

func (bc *GitIgnores) SetConfig(rawConfig interface{}) error {
	var config Config
	err := mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	bc.config = &config

	return nil
}

func (bc *GitIgnores) GetConfig() rule.Config {
	return bc.config
}

func (bc *GitIgnores) Run(ctx *rule.Context) error {

	fBytes, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return err
	}

	ignoredEntries := strings.Split(string(fBytes), "\n")

	var found bool
	for _, expectedEntry := range bc.config.Entries {
		found = false
		for _, ignoredEntry := range ignoredEntries {
			if ignoredEntry == expectedEntry {
				found = true
				break
			}
		}

		if !found {
			ctx.AddFailure(rule.NewBasicFailure(
				bc,
				fmt.Sprintf("expected git to ignore: '%s'", expectedEntry),
			), bc.config.IsWarn())
		}
	}

	return nil
}
