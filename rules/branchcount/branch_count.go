package branchcount

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"strconv"
)

const (
	ID   = "GL001"
	Name = "branch_count"
)

type BranchCount struct {
	Repo   *git.Repository
	config *Config
}

func (bc *BranchCount) ID() string {
	return ID
}

func (bc *BranchCount) Name() string {
	return Name
}

func (bc *BranchCount) SetConfig(rawConfig interface{}) error {
	var config Config
	err := mapstructure.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	bc.config = &config

	return nil
}

func (bc *BranchCount) GetConfig() rule.Config {
	return bc.config
}

func (bc *BranchCount) Run(ctx *rule.Context) error {
	branches, err := bc.Repo.Branches(&git.BranchOpts{
		ShortNames: true,
	})
	if err != nil {
		return err
	}

	bCount := len(branches)
	if bCount > bc.config.Max {
		ctx.AddFailure(rule.NewActualExpectedFailure(
			bc,
			"greater than permitted value",
			strconv.Itoa(bCount),
			fmt.Sprintf("less than %d", bc.config.Max),
		), bc.config.IsWarn())
	}

	return nil
}
