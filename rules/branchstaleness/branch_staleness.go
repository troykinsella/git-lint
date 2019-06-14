package branchstaleness

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"github.com/troykinsella/sloppy-duration"
	"time"
)

const (
	ID = "GL003"
	Name = "branch_staleness"
)

type BranchStaleness struct {
	Repo *git.Repository
	config *Config
}

func (bs *BranchStaleness) ID() string {
	return ID
}

func (bs *BranchStaleness) Name() string {
	return Name
}

func (bs *BranchStaleness) SetConfig(rawConfig interface{}) error {

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

	bs.config = &config
	return nil
}

func (bs *BranchStaleness) GetConfig() rule.Config {
	return bs.config
}

func (bs *BranchStaleness) Run(ctx *rule.Context) error {
	branches, err := bs.Repo.Branches(&git.BranchOpts{
		ShortNames: false,
	})
	if err != nil {
		return err
	}

	now := time.Now()

	for _, branch := range branches {
		lastCommitTime, err := bs.Repo.LatestCommitTime(branch)
		if err != nil {
			return err
		}

		actualStaleness := now.Sub(lastCommitTime)
		if actualStaleness > bs.config.maxDuration.Duration() {
			ctx.AddFailure(rule.NewActualExpectedFailure(
				bs,
				fmt.Sprintf("last commit to '%s' older than permitted staleness duration", bs.Repo.ShortName(branch)),
				sloppy_duration.Wrap(actualStaleness),
				fmt.Sprintf("less than %s", bs.config.maxDuration.String()),
			), bs.config.IsWarn())
		}
	}

	return nil
}
