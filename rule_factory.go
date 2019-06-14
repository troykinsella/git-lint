package git_lint

import (
	"fmt"
	"github.com/troykinsella/git-lint/git"
	"github.com/troykinsella/git-lint/rule"
	"github.com/troykinsella/git-lint/rules/branchcount"
	"github.com/troykinsella/git-lint/rules/branchlastcommit"
	"github.com/troykinsella/git-lint/rules/branchname"
	"github.com/troykinsella/git-lint/rules/branchsingleton"
	"github.com/troykinsella/git-lint/rules/gitignore"
	"github.com/troykinsella/git-lint/rules/gitkeep"
	"github.com/troykinsella/git-lint/rules/tagname"
)

func NewRule(name string, repo *git.Repository) rule.Rule {
	var r rule.Rule

	switch name {
	case branchcount.Name:
		r = &branchcount.BranchCount{Repo: repo}
	case branchlastcommit.Name:
		r = &branchlastcommit.BranchLastCommit{Repo: repo}
	case branchname.Name:
		r = &branchname.BranchName{Repo: repo}
	case branchsingleton.Name:
		r = &branchsingleton.BranchSingleton{Repo: repo}
	case gitignore.Name:
		r = &gitignore.GitIgnore{}
	case gitkeep.Name:
		r = &gitkeep.GitKeep{}
	case tagname.Name:
		r = &tagname.TagName{Repo: repo}
	default:
		panic(fmt.Errorf("unknown rule name: %s", name))
	}

	return r
}
