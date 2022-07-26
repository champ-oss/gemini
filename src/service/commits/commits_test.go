package commits

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

// Test_ProcessRepo_Success checks that ProcessRepo() is successful
func Test_ProcessRepo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	branch := "branch1"
	commits := []*model.Commit{
		{
			Owner: owner,
		},
	}

	gitClient.EXPECT().GetCommits(owner, name, branch).Return(commits, nil)
	repo.EXPECT().AddCommits(commits).Return(int64(1), nil)
	svc := service{repo: repo, gitClient: gitClient}
	assert.Nil(t, svc.ProcessRepo(owner, name, branch))
}

// Test_ProcessRepo_NeverSynced_WithError tests ProcessRepo() with no previous sync time and an error result
func Test_ProcessRepo_NeverSynced_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	branch := "branch1"

	gitClient.EXPECT().GetCommits(owner, name, branch).Return(nil, fmt.Errorf("test error"))
	svc := service{repo: repo, gitClient: gitClient}
	assert.Error(t, svc.ProcessRepo(owner, name, branch))
}
