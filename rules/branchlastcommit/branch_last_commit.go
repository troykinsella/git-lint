package branchlastcommit

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"github.com/troykinsella/sloppy-duration"
	"time"
)

const (
	ID   = "GL002"
	Name = "branch_last_commit"
)

type BranchLastCommit struct {
	Repo   *git.Repository
	config *Config
}

func (blc *BranchLastCommit) ID() string {
	return ID
}

func (blc *BranchLastCommit) Name() string {
	return Name
}

func (blc *BranchLastCommit) SetConfig(rawConfig interface{}) error {

	var config Config
	var err error

	err = mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	if config.Max_Duration == "" {
		return rule.NewConfigError(Name + ": max_duration must be supplied")
	}

	config.maxDuration, err = sloppy_duration.Parse(config.Max_Duration)
	if err != nil {
		return rule.NewConfigError(Name + ": invalid duration: '" + config.Max_Duration + "'")
	}

	blc.config = &config
	return nil
}

func (blc *BranchLastCommit) GetConfig() rule.Config {
	return blc.config
}

func (blc *BranchLastCommit) Run(ctx *rule.Context) error {
	branches, err := blc.Repo.Branches(&git.BranchOpts{
		ShortNames: false,
	})
	if err != nil {
		return err
	}

	now := time.Now()

	for _, branch := range branches {
		lastCommitTime, err := blc.Repo.LatestCommitTime(branch)
		if err != nil {
			return err
		}

		actualStaleness := now.Sub(lastCommitTime)
		if actualStaleness > blc.config.maxDuration.Duration() {
			ctx.AddFailure(rule.NewActualExpectedFailure(
				blc,
				fmt.Sprintf("last commit to '%s' older than permitted duration", blc.Repo.ShortName(branch)),
				sloppy_duration.Wrap(actualStaleness),
				fmt.Sprintf("less than %s", blc.config.maxDuration.String()),
			), blc.config.IsWarn())
		}
	}

	return nil
}
