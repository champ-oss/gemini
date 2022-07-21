package actions

import (
	"fmt"
	"github.com/champ-oss/gemini/mocks/mock_adapter"
	"github.com/champ-oss/gemini/mocks/mock_repository"
	"github.com/champ-oss/gemini/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	svc := NewService(nil, nil)
	assert.NotNil(t, svc)
}

// Test_ProcessRepo_Success checks that processRepo() is successful
func Test_service_ProcessRepo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	runs := []*model.WorkflowRun{
		{
			Owner:  owner,
			Repo:   name,
			Branch: branch,
		},
	}

	gitClient.EXPECT().GetWorkflowRuns(owner, name, branch).Return(runs, nil)
	repo.EXPECT().AddWorkflowRuns(runs).Return(int64(1), nil)

	svc := service{repo: repo, gitClient: gitClient}
	assert.Nil(t, svc.ProcessRepo(owner, name, branch))
}

// Test_service_ProcessRepo_Reruns checks that processRepo() is successful with reruns
func Test_service_ProcessRepo_Reruns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	runs := []*model.WorkflowRun{
		{
			Owner:      owner,
			Repo:       name,
			Branch:     branch,
			Name:       "workflow1",
			RunAttempt: 1,
			RunID:      111,
		},
		{
			Owner:      owner,
			Repo:       name,
			Branch:     branch,
			Name:       "workflow2",
			RunAttempt: 2,
			RunID:      222,
		},
	}

	// Expected rerun
	reruns := []*model.WorkflowRun{
		{
			Owner:      owner,
			Repo:       name,
			Branch:     branch,
			Name:       "workflow2",
			RunAttempt: 1,
			RunID:      222,
		},
	}

	gitClient.EXPECT().GetWorkflowRuns(owner, name, branch).Return(runs, nil)
	repo.EXPECT().AddWorkflowRuns(runs).Return(int64(1), nil)
	gitClient.EXPECT().GetWorkflowRunAttempt(owner, name, branch, int64(222), 1).Return(reruns[0], nil)
	repo.EXPECT().AddWorkflowRuns(reruns).Return(int64(1), nil)

	svc := service{repo: repo, gitClient: gitClient}
	assert.Nil(t, svc.ProcessRepo(owner, name, branch))
}

// Test_ProcessRepo_Error checks that processRepo() is fails
func Test_service_ProcessRepo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	branch := "branch1"

	gitClient.EXPECT().GetWorkflowRuns(owner, name, branch).Return(nil, fmt.Errorf("test error"))

	svc := service{repo: repo, gitClient: gitClient}
	assert.Error(t, svc.ProcessRepo(owner, name, branch))
}

// Test_service_getRerunsForWorkflow_Success tests getRerunsForWorkflow successfully
func Test_service_getRerunsForWorkflow_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	svc := service{repo: repo, gitClient: gitClient}

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	run := &model.WorkflowRun{
		Owner:      owner,
		Repo:       name,
		Branch:     branch,
		Name:       "workflow1",
		RunAttempt: 3,
		RunID:      198675491,
	}

	gitClient.EXPECT().GetWorkflowRunAttempt(owner, name, branch, run.RunID, run.RunAttempt-1)
	gitClient.EXPECT().GetWorkflowRunAttempt(owner, name, branch, run.RunID, run.RunAttempt-2)
	results, err := svc.getRerunsForWorkflow(owner, name, branch, run)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

// Test_service_getRerunsForWorkflow_None tests getRerunsForWorkflow when no reruns are present
func Test_service_getRerunsForWorkflow_None(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	svc := service{repo: repo, gitClient: gitClient}

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	run := &model.WorkflowRun{
		Owner:      owner,
		Repo:       name,
		Branch:     branch,
		Name:       "workflow1",
		RunAttempt: 1,
		RunID:      198675491,
	}

	results, err := svc.getRerunsForWorkflow(owner, name, branch, run)
	assert.NoError(t, err)
	assert.Len(t, results, 0)
}

// Test_service_getRerunsForWorkflow_Error tests getRerunsForWorkflow with an error response
func Test_service_getRerunsForWorkflow_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	svc := service{repo: repo, gitClient: gitClient}

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	run := &model.WorkflowRun{
		Owner:      owner,
		Repo:       name,
		Branch:     branch,
		Name:       "workflow1",
		RunAttempt: 2,
		RunID:      198675491,
	}

	gitClient.EXPECT().GetWorkflowRunAttempt(owner, name, branch, run.RunID, run.RunAttempt-1).Return(nil, fmt.Errorf("test error"))
	results, err := svc.getRerunsForWorkflow(owner, name, branch, run)
	assert.Error(t, err)
	assert.Nil(t, results)
}

// Test_service_getAllWorkflowReruns tests that getting reruns are successful
func Test_service_getAllWorkflowReruns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)
	svc := service{repo: repo, gitClient: gitClient}

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	runs := []*model.WorkflowRun{
		{
			Owner:      owner,
			Repo:       name,
			Branch:     branch,
			Name:       "workflow1",
			RunAttempt: 1,
			RunID:      111,
		},
		{
			Owner:      owner,
			Repo:       name,
			Branch:     branch,
			Name:       "workflow2",
			RunAttempt: 2,
			RunID:      222,
		},
	}

	// Expected rerun
	reruns := []*model.WorkflowRun{
		{
			Owner:      owner,
			Repo:       name,
			Branch:     branch,
			Name:       "workflow2",
			RunAttempt: 1,
			RunID:      222,
		},
	}

	gitClient.EXPECT().GetWorkflowRunAttempt(owner, name, branch, int64(222), 1).Return(reruns[0], nil)
	results, err := svc.getAllWorkflowReruns(owner, name, branch, runs)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, reruns, results)
}
