package common

import (
	"github.com/champ-oss/gemini/mocks/mock_adapter"
	"github.com/champ-oss/gemini/mocks/mock_repository"
	"github.com/champ-oss/gemini/mocks/mock_service_actions"
	"github.com/champ-oss/gemini/mocks/mock_service_commits"
	"github.com/champ-oss/gemini/mocks/mock_service_pullrequests"
	"github.com/champ-oss/gemini/mocks/mock_service_terraformrefs"
	"github.com/champ-oss/gemini/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	svc := NewService(nil, nil, 1)
	assert.NotNil(t, svc)
}

// Test_PopulateDefaultBranch_Success tests that populateDefaultBranch() adds the branch to the repo
func Test_PopulateDefaultBranch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	svc := service{repo: repo, gitClient: gitClient}

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	repos := []*model.Repo{
		{
			Owner: owner,
			Name:  name,
		},
	}

	gitClient.EXPECT().GetDefaultBranch(owner, name).Return(branch, nil)
	assert.Nil(t, svc.populateDefaultBranch(repos))
	assert.Equal(t, repos[0].Branch, branch)
}

// Test_PopulateDefaultBranch_NotEmpty tests that populateDefaultBranch() doesn't run when the branch is not empty
func Test_PopulateDefaultBranch_NotEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	svc := service{repo: repo, gitClient: gitClient}

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	repos := []*model.Repo{
		{
			Owner:  owner,
			Name:   name,
			Branch: branch,
		},
	}
	assert.Nil(t, svc.populateDefaultBranch(repos))
	assert.Equal(t, repos[0].Branch, branch)
}

func Test_StartGemini_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	commitsService := mock_service_commits.NewMockServiceInterface(ctrl)
	actionsService := mock_service_actions.NewMockServiceInterface(ctrl)
	terraformrefsService := mock_service_terraformrefs.NewMockServiceInterface(ctrl)
	pullrequestService := mock_service_pullrequests.NewMockServiceInterface(ctrl)
	svc := service{
		repo:                 repo,
		gitClient:            gitClient,
		commitsService:       commitsService,
		actionsService:       actionsService,
		terraformrefsService: terraformrefsService,
		pullRequests:         pullrequestService,
	}
	svc.runOnce = true

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	repos := []*model.Repo{
		{
			Owner: owner,
			Name:  name,
		},
	}

	gitClient.EXPECT().GetDefaultBranch(owner, name).Return(branch, nil)
	commitsService.EXPECT().ProcessRepo(owner, name, branch).Return(nil)
	actionsService.EXPECT().ProcessRepo(owner, name, branch).Return(nil)
	terraformrefsService.EXPECT().ProcessRepo(owner, name, branch).Return(nil)
	pullrequestService.EXPECT().ProcessRepo(owner, name).Return(nil)

	assert.Nil(t, svc.StartGemini(repos))
}
