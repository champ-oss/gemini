package actions

import (
	"github.com/champ-oss/gemini/adapter"
	"github.com/champ-oss/gemini/model"
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
	log.Info("Syncing workflow runs to the database")
	runs, err := s.gitClient.GetWorkflowRuns(owner, name, branch)
	if err != nil {
		return err
	}

	if len(runs) > 0 {
		_, err := s.repo.AddWorkflowRuns(runs)
		if err != nil {
			return err
		}
		log.Infof("Synced %d workflow runs to the database", len(runs))
	}

	log.Info("Checking for workflow re-runs to sync to the database")
	reruns, err := s.getAllWorkflowReruns(owner, name, branch, runs)
	if err != nil {
		return err
	}
	if len(reruns) > 0 {
		_, err = s.repo.AddWorkflowRuns(reruns)
		if err != nil {
			return err
		}
		log.Infof("Synced %d workflow re-runs to the database", len(reruns))
	}
	return nil
}

// getAllWorkflowReruns checks for and fetches reruns for each workflow run
func (s *service) getAllWorkflowReruns(owner string, name string, branch string, runs []*model.WorkflowRun) ([]*model.WorkflowRun, error) {
	var allReruns []*model.WorkflowRun

	for _, run := range runs {
		if run.RunAttempt > 1 {
			runs, err := s.getRerunsForWorkflow(owner, name, branch, run)
			if err != nil {
				return nil, err
			}
			allReruns = append(allReruns, runs...)
		}
	}

	return allReruns, nil
}

// getRerunsForWorkflow gets all previous workflow reruns for the given workflow run
func (s *service) getRerunsForWorkflow(owner string, name string, branch string, run *model.WorkflowRun) ([]*model.WorkflowRun, error) {
	var reruns []*model.WorkflowRun
	rerunAttempt := run.RunAttempt

	for {
		rerunAttempt--
		if rerunAttempt < 1 {
			break
		}

		log.Infof("Getting attempt #%d for run %d, workflow: %s", rerunAttempt, run.RunID, run.Name)
		rerun, err := s.gitClient.GetWorkflowRunAttempt(owner, name, branch, run.RunID, rerunAttempt)
		if err != nil {
			return nil, err
		}

		reruns = append(reruns, rerun)
	}
	return reruns, nil
}
