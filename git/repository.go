package git

import (
	gitlib "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"strings"
	"time"
)

type BranchOpts struct {
	ShortNames bool
}

type Repository struct {
	repo   *gitlib.Repository
	remote string
}

func NewRepository(path string, remote string) (*Repository, error) {

	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	repo, err := gitlib.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &Repository{
		repo:   repo,
		remote: remote,
	}, nil
}

func (r *Repository) Remote() string {
	return r.remote
}

func (r *Repository) refsForPrefix(prefix string, shortNames bool) ([]string, error) {

	refItr, err := r.repo.References()
	if err != nil {
		return nil, err
	}

	refs := make([]string, 0)

	err = refItr.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			name := ref.Name().String()
			if strings.HasPrefix(name, prefix) {
				if shortNames {
					name = name[len(prefix)+1:]
				}
				refs = append(refs, name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return refs, nil
}

func (r *Repository) Branches(opts *BranchOpts) ([]string, error) {
	prefix := "refs/remotes/" + r.remote
	return r.refsForPrefix(prefix, opts.ShortNames)
}

func (r *Repository) Tags() ([]string, error) {
	prefix := "refs/tags"
	return r.refsForPrefix(prefix, true)
}

func (r *Repository) ShortName(refName string) string {
	prefix := "refs/remotes/" + r.remote
	if strings.HasPrefix(refName, prefix) {
		return refName[len(prefix)+1:]
	}

	prefix = "refs/tags"
	if strings.HasPrefix(refName, prefix) {
		return refName[len(prefix)+1:]
	}

	return refName
}

func (r *Repository) FirstCommitTime(refName string) (time.Time, error) {

	ref, err := r.repo.Reference(plumbing.ReferenceName(refName), false)
	if err != nil {
		return time.Time{}, err
	}

	comItr, err := r.repo.Log(&gitlib.LogOptions{
		From:  ref.Hash(),
		Order: gitlib.LogOrderCommitterTime,
	})
	if err != nil {
		return time.Time{}, err
	}

	var commit *object.Commit
	for {
		c, err := comItr.Next()
		if err != nil {
			return time.Time{}, err
		}

		if c == nil {
			break
		} else {
			commit = c
		}
	}

	return commit.Committer.When, nil
}

func (r *Repository) LatestCommitTime(refName string) (time.Time, error) {

	ref, err := r.repo.Reference(plumbing.ReferenceName(refName), false)
	if err != nil {
		return time.Time{}, err
	}

	comItr, err := r.repo.Log(&gitlib.LogOptions{
		From:  ref.Hash(),
		Order: gitlib.LogOrderCommitterTime,
	})
	if err != nil {
		return time.Time{}, err
	}

	commit, err := comItr.Next()
	if err != nil {
		return time.Time{}, err
	}

	return commit.Committer.When, nil
}
