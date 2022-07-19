package pullrequests

import (
	"fmt"
	"github.com/champ-oss/gemini/mocks/mock_adapter"
	"github.com/champ-oss/gemini/mocks/mock_repository"
	"github.com/champ-oss/gemini/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCommitModel = &model.Commit{
	Owner:          "owner1",
	Repo:           "repo1",
	Branch:         "branch1",
	Message:        "commit msg",
	CommitterName:  "test committer",
	CommitterDate:  1642708819,
	CommitterEmail: "testcommitter@test.com",
	AuthorName:     "test author",
	AuthorDate:     1642708819,
	AuthorEmail:    "testauthor@test.com",
	Url:            "http://localhost",
	Sha:            "702fe8ee76f422edd4bc257a0a2171af26563063",
}

var testPullRequestModel = &model.PullRequest{
	Owner:                "owner1",
	Repo:                 "repo1",
	MergeCommitSHA:       "702fe8ee76f422edd4bc257a0a2171af26563063",
	Number:               10,
	State:                "open",
	Title:                "test pull request",
	PullRequestCreatedAt: 1642708819,
	PullRequestUpdatedAt: 1642708819,
	PullRequestClosedAt:  1642708819,
	PullRequestMergedAt:  1642708819,
	Draft:                false,
	Merged:               false,
	Commits:              3,
}

var testPullRequestCommitModel = &model.PullRequestCommit{
	Owner:          "owner1",
	Repo:           "repo1",
	MergeCommitSHA: "702fe8ee76f422edd4bc257a0a2171af26563063",
	Number:         10,
	Sha:            "702fe8ee76f422edd4bc257a0a2171af26563063",
	CommitterDate:  1642708819,
}

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

	gitClient.EXPECT().GetPullRequests(owner, name, "closed").Return([]*model.PullRequest{testPullRequestModel}, nil)
	gitClient.EXPECT().GetPullRequestCommits(owner, name, 10).Return([]*model.Commit{testCommitModel}, nil)
	repo.EXPECT().AddPullRequestCommits([]*model.PullRequestCommit{testPullRequestCommitModel}).Return(int64(1), nil)

	svc := service{repo: repo, gitClient: gitClient}
	assert.Nil(t, svc.ProcessRepo(owner, name))
}

// Test_ProcessRepo_Success checks that processRepo() is successful
func Test_service_ProcessRepo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repository.NewMockRepository(ctrl)
	gitClient := mock_adapter.NewMockGitClient(ctrl)

	owner := "owner1"
	name := "repo1"
	svc := service{repo: repo, gitClient: gitClient}

	gitClient.EXPECT().GetPullRequests(owner, name, "closed").Return(nil, fmt.Errorf("test error"))
	assert.Error(t, svc.ProcessRepo(owner, name))

	gitClient.EXPECT().GetPullRequests(owner, name, "closed").Return([]*model.PullRequest{testPullRequestModel}, nil)
	gitClient.EXPECT().GetPullRequestCommits(owner, name, 10).Return(nil, fmt.Errorf("test error"))
	assert.Error(t, svc.ProcessRepo(owner, name))
}

// Test_parsePullRequestCommits tests parsing pull request commits into the PullRequestCommit model
func Test_parsePullRequestCommits(t *testing.T) {
	parsed := parsePullRequestCommits("owner1", "repo1", testPullRequestModel, []*model.Commit{testCommitModel})
	assert.Equal(t, []*model.PullRequestCommit{testPullRequestCommitModel}, parsed)
}
