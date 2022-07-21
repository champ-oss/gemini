package common

import (
	"github.com/champ-oss/gemini/adapter"
	"github.com/champ-oss/gemini/model"
	"github.com/champ-oss/gemini/repository"
	"github.com/champ-oss/gemini/service/actions"
	"github.com/champ-oss/gemini/service/commits"
	"github.com/champ-oss/gemini/service/pullrequests"
	"github.com/champ-oss/gemini/service/terraformrefs"
	log "github.com/sirupsen/logrus"
	"time"
)

type service struct {
	repo                 repository.Repository
	gitClient            adapter.GitClient
	commitsService       commits.ServiceInterface
	actionsService       actions.ServiceInterface
	terraformrefsService terraformrefs.ServiceInterface
	pullRequests         pullrequests.ServiceInterface
	minutesBetweenChecks float64
	runOnce              bool
}

// NewService initializes a new service
func NewService(repo repository.Repository, gitClient adapter.GitClient, minutesBetweenChecks float64) *service {
	commitsService := commits.NewService(repo, gitClient)
	actionsService := actions.NewService(repo, gitClient)
	terraformrefsService := terraformrefs.NewService(repo, gitClient)
	pullRequests := pullrequests.NewService(repo, gitClient)

	return &service{
		repo:                 repo,
		gitClient:            gitClient,
		commitsService:       commitsService,
		actionsService:       actionsService,
		terraformrefsService: terraformrefsService,
		pullRequests:         pullRequests,
		minutesBetweenChecks: minutesBetweenChecks,
	}
}

// StartGemini is the entrypoint for the overall Gemini service
func (s *service) StartGemini(repos []*model.Repo) error {
	for {
		err := s.populateDefaultBranch(repos)
		if err != nil {
			return err
		}

		for _, repo := range repos {

			log.Infof("Starting to process repo: %s/%s (%s branch)", repo.Owner, repo.Name, repo.Branch)
			err := s.commitsService.ProcessRepo(repo.Owner, repo.Name, repo.Branch)
			if err != nil {
				return err
			}

			err = s.actionsService.ProcessRepo(repo.Owner, repo.Name, repo.Branch)
			if err != nil {
				return err
			}

			err = s.terraformrefsService.ProcessRepo(repo.Owner, repo.Name, repo.Branch)
			if err != nil {
				return err
			}

			err = s.pullRequests.ProcessRepo(repo.Owner, repo.Name)
			if err != nil {
				return err
			}

			log.Infof("Done processing: %s", repo.Name)
		}

		log.Infof("Done processing all repos")
		if s.runOnce {
			log.Info("Exiting since runOnce is set to true")
			return nil
		}
		log.Infof("Sleeping until next check in %f minutes", s.minutesBetweenChecks)
		time.Sleep(time.Minute * time.Duration(s.minutesBetweenChecks))
	}
}

// populateDefaultBranch populates the branch field for each repo with the default branch
func (s *service) populateDefaultBranch(repos []*model.Repo) error {
	for _, repo := range repos {
		if repo.Branch != "" {
			continue
		}
		log.Infof("Querying for default branch for repo: %s/%s", repo.Owner, repo.Name)
		branch, err := s.gitClient.GetDefaultBranch(repo.Owner, repo.Name)
		if err != nil {
			return err
		}
		repo.Branch = branch
	}
	return nil
}
