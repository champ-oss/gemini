package pullrequests

import (
	"github.com/champ-oss/gemini/adapter"
	"github.com/champ-oss/gemini/model"
	"github.com/champ-oss/gemini/repository"
	log "github.com/sirupsen/logrus"
)

type ServiceInterface interface {
	ProcessRepo(owner string, name string) error
}

type service struct {
	repo      repository.Repository
	gitClient adapter.GitClient
}

// NewService initializes a new service
func NewService(repo repository.Repository, gitClient adapter.GitClient) *service {
	return &service{
		repo:      repo,
		gitClient: gitClient,
	}
}

// ProcessRepo process one repo
func (s *service) ProcessRepo(owner string, name string) error {
	log.Info("Syncing pull requests to the database")
	pullRequests, err := s.gitClient.GetPullRequests(owner, name, "closed")
	if err != nil {
		return err
	}

	for _, pr := range pullRequests {
		log.Infof("Syncing commits for pull request: %d", pr.Number)
		commits, err := s.gitClient.GetPullRequestCommits(owner, name, pr.Number)
		if err != nil {
			return err
		}
		pullRequestCommits := parsePullRequestCommits(owner, name, pr, commits)
		_, err = s.repo.AddPullRequestCommits(pullRequestCommits)
		if err != nil {
			return err
		}
		log.Infof("Synced %d pull request commits to the database", len(commits))
	}

	return nil
}

// parsePullRequestCommits parses commits into the PullRequestCommit model
func parsePullRequestCommits(owner string, name string, pr *model.PullRequest, commits []*model.Commit) []*model.PullRequestCommit {
	var parsed []*model.PullRequestCommit
	for _, commit := range commits {
		parsed = append(parsed, &model.PullRequestCommit{
			Owner:          owner,
			Repo:           name,
			Number:         pr.Number,
			MergeCommitSHA: pr.MergeCommitSHA,
			Sha:            commit.Sha,
			CommitterDate:  commit.CommitterDate,
		})
	}
	return parsed
}
