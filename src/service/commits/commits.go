package commits

import (
	"github.com/champ-oss/gemini/adapter"
	"github.com/champ-oss/gemini/repository"
	log "github.com/sirupsen/logrus"
)

type ServiceInterface interface {
	ProcessRepo(owner string, name string, branch string) error
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
func (s *service) ProcessRepo(owner string, name string, branch string) error {
	log.Info("Syncing commits to the database")
	commits, err := s.gitClient.GetCommits(owner, name, branch)
	if err != nil {
		return err
	}

	if len(commits) > 0 {
		_, err := s.repo.AddCommits(commits)
		if err != nil {
			return err
		}
		log.Infof("Synced %d commits to the database", len(commits))
	}
	return nil
}
