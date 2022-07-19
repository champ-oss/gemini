package terraformrefs

import (
	"github.com/champ-oss/gemini/adapter"
	"github.com/champ-oss/gemini/repository"
	log "github.com/sirupsen/logrus"
)

const terraformFile = "main.tf"

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
	log.Info("Getting terraform apply workflow runs")
	applyRuns := s.repo.GetWorkflowRunsByName(owner, name, branch, "apply")
	log.Infof("Retrieved %d terraform apply workflow runs", len(applyRuns))

	for _, applyRun := range applyRuns {
		log.Infof("Getting %s contents for commit: %s", terraformFile, applyRun.Sha)
		contents, err := s.gitClient.GetContents(owner, name, terraformFile, applyRun.Sha)
		if err != nil {
			log.Warnf("Failed to get contents of %s for commit: %s - %s", terraformFile, applyRun.Sha, err)
			continue
		}

		hcl, err := parseStringAsHcl(contents, terraformFile)
		if err != nil {
			log.Warnf("Failed to parse string as HCL: %s. Contents: \n%s", err, contents)
			continue
		}

		// Parse the Terraform file and get a list of module call ref information
		refs, err := parseModuleCalls(hcl, applyRun)
		if err != nil {
			return err
		}

		if len(refs) > 0 {
			_, err := s.repo.AddTerraformRefs(refs)
			if err != nil {
				return err
			}
			log.Infof("Synced %d terraform refs to the database", len(refs))
		}
	}
	return nil
}
