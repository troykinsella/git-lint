package gitkeep

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/rule"
	"io/ioutil"
	"os"
)

const (
	ID   = "GL006"
	Name = "git_keep"
)

type GitKeep struct {
	config *Config
}

func (gk *GitKeep) ID() string {
	return ID
}

func (gk *GitKeep) Name() string {
	return Name
}

func (gk *GitKeep) SetConfig(rawConfig interface{}) error {
	var config Config
	err := mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	gk.config = &config

	return nil
}

func (gk *GitKeep) GetConfig() rule.Config {
	return gk.config
}

func (gk *GitKeep) Run(ctx *rule.Context) error {
	for _, dir := range gk.config.Directories {
		cool, err := isPopulatedDir(dir)
		if err != nil {
			return err
		}

		if !cool {
			ctx.AddFailure(rule.NewBasicFailure(
				gk,
				fmt.Sprintf("expected directory to exist and contain at least one file: %s", dir),
			), gk.config.IsWarn())
		}
	}

	return nil
}

func isPopulatedDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	if !stat.IsDir() {
		return false, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	if len(files) == 0 {
		return false, nil
	}

	return true, nil
}
