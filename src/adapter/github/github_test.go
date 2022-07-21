package github

import (
	"github.com/champ-oss/gemini/model"
	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testUnixTime = int64(1642708819)
var testTime = time.Unix(testUnixTime, 0)
var owner = "owner1"
var repo = "repo1"
var branch = "branch1"
var sha = "702fe8ee76f422edd4bc257a0a2171af26563063"

var testCommitModel = &model.Commit{
	Owner:          owner,
	Repo:           repo,
	Branch:         branch,
	Message:        "commit msg",
	CommitterName:  "test committer",
	CommitterDate:  testUnixTime,
	CommitterEmail: "testcommitter@test.com",
	AuthorName:     "test author",
	AuthorDate:     testUnixTime,
	AuthorEmail:    "testauthor@test.com",
	Url:            "http://localhost",
	Sha:            sha,
}

var testCommitGitHub = &github.RepositoryCommit{
	SHA: github.String(testCommitModel.Sha),
	URL: github.String(testCommitModel.Url),
	Commit: &github.Commit{
		Message: github.String(testCommitModel.Message),
		Committer: &github.CommitAuthor{
			Name:  github.String(testCommitModel.CommitterName),
			Date:  &testTime,
			Email: github.String(testCommitModel.CommitterEmail),
		},
		Author: &github.CommitAuthor{
			Name:  github.String(testCommitModel.AuthorName),
			Date:  &testTime,
			Email: github.String(testCommitModel.AuthorEmail),
		},
	},
}

var testWorkflowRun = &github.WorkflowRun{
	Name:       github.String("test run"),
	HeadSHA:    github.String(sha),
	HeadBranch: github.String(branch),
	Conclusion: github.String("success"),
	CreatedAt:  &github.Timestamp{Time: testTime},
	UpdatedAt:  &github.Timestamp{Time: testTime},
}

var testWorkflowRunModel = &model.WorkflowRun{
	Owner:        owner,
	Repo:         repo,
	Name:         "test run",
	Branch:       branch,
	Sha:          sha,
	Conclusion:   "success",
	RunCreatedAt: testUnixTime,
	RunUpdatedAt: testUnixTime,
}

var testPullRequest = &github.PullRequest{
	MergeCommitSHA: github.String(sha),
	Number:         github.Int(10),
	State:          github.String("open"),
	Title:          github.String("test pull request"),
	CreatedAt:      &testTime,
	UpdatedAt:      &testTime,
	ClosedAt:       &testTime,
	MergedAt:       &testTime,
	Draft:          github.Bool(false),
	Merged:         github.Bool(false),
	Commits:        github.Int(3),
}

var testPullRequestModel = &model.PullRequest{
	Owner:                owner,
	Repo:                 repo,
	MergeCommitSHA:       sha,
	Number:               10,
	State:                "open",
	Title:                "test pull request",
	PullRequestCreatedAt: testUnixTime,
	PullRequestUpdatedAt: testUnixTime,
	PullRequestClosedAt:  testUnixTime,
	PullRequestMergedAt:  testUnixTime,
	Draft:                false,
	Merged:               false,
	Commits:              3,
}

// Needs to be added to upstream github.com/migueleliasweb/go-github-mock
var GetWorkflowRunAttempt mock.EndpointPattern = mock.EndpointPattern{
	Pattern: "/repos/{owner}/{repo}/actions/runs/{run_id}/attempts/{attempt_number}",
	Method:  "GET",
}

// Test_NewAdapter tests creating a new adapter by passing a dummy private key
func Test_NewAdapter(t *testing.T) {
	adapter, err := NewAdapter(9999, 9999, []byte("-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQC2UrmMz+fXN2FCZLCATJJEhJrYVZyXmvYcnMwr3VLN+QxM8WKJ\nOk0HSixah4Iw45nYysMmFouGRvukVCE+fm0oH4fEcDVt9MH6+lteeEH/1X5rti6b\nUkWw0DrKfGOmyXTY5BXxRUd5czdig/I8f2rh4i4+urWCGzpDNUutmfp9JwIDAQAB\nAoGAbjM6K75OZ2r1wmeRtzqQ/hEYdsHoUEo9j7XGQo0Xy59IrAkKgd9XR5yxilZ1\nfo9TIhIMOi1OT+7/kqe3IErSNnUNqzoC77AypYnDFl5uV6RxyJynvYQXlsBTE2LJ\nT4JicDxVXnbJBYJK2ioAdqN9vug7aOsk1SymZaDKf0F7wuECQQDYRvcLDEF+hYCA\nWay9Pk1uT6V2wRSqb7q23eZiWjhrOWFwBQJDMoCi6vvFBzVo3ICYUXucXmG2U0f1\n1fGERc/LAkEA189KNGRwP4Vdx5Snqi5WVJzXD5Nov1O09Ks4fgyEtIJvk6vWz5Ko\nCNTqjaEq/r7B1Aul5AJ9NvFBXEG/oGAklQJBANT+XoF02nNdysWciu/8kYkXyx5+\n3HlVe45oTmGB9Jo0cm89n5LKA8FupfDOPp08uzBG3vOKR7Slo/LJdgcMMa0CQQCz\n77bsNi5NGDLX/H9LarU6eUbrSroUhIOlWLmSh3eCVhsX4jgJ/Dq0mmoyyoVhv8U2\nuruHf/fM/pzDgmJ3IpJ9AkBDj4PL9gUaxwEjcG+J3InbRZyFSCETugj9eDJVxZk4\nlV8gKSiObUDUb2BP3oRzFOLqYmxoxsgMKxEcqleSG851\n-----END RSA PRIVATE KEY-----\n"))
	assert.NoError(t, err)
	assert.NotNil(t, adapter.client)
}

// Test_NewAdapter_InvalidPem tests creating a new adapter by passing an invalid private key
func Test_NewAdapter_InvalidPem(t *testing.T) {
	adapter, err := NewAdapter(9999, 9999, []byte("nonsense"))
	assert.Error(t, err)
	assert.Nil(t, adapter)
}

// Test_NewAdapter_NoCredentials tests creating a new adapter by not passing any credentials
func Test_NewAdapter_NoCredentials(t *testing.T) {
	adapter, err := NewAdapter(0, 0, []byte(""))
	assert.NoError(t, err)
	assert.NotNil(t, adapter.client)
}

// Test_GetCommits tests getting all the test commits
func Test_GetCommits(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchPages(
			mock.GetReposCommitsByOwnerByRepo,
			[]*github.RepositoryCommit{testCommitGitHub, testCommitGitHub},
			[]*github.RepositoryCommit{testCommitGitHub, testCommitGitHub},
		),
	)

	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	results, err := gh.GetCommits("owner1", "repo1", "branch1")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(results))
}

// Test_ParseGithubCommits tests parsing testCommitGitHub into testCommitModel
func Test_ParseGithubCommits(t *testing.T) {
	parsed := parseGithubCommits(testCommitModel.Owner, testCommitModel.Repo,
		testCommitModel.Branch, []*github.RepositoryCommit{testCommitGitHub})
	assert.Equal(t, []*model.Commit{testCommitModel}, parsed)
}

// Test_GetDefaultBranch tests getting the default branch for a repo
func Test_GetDefaultBranch(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposByOwnerByRepo,
			&github.Repository{
				DefaultBranch: github.String("test"),
			},
		))
	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	result, err := gh.GetDefaultBranch(testCommitModel.Owner, testCommitModel.Repo)
	assert.Nil(t, err)
	assert.Equal(t, "test", result)
}

// Test_GetDefaultBranch_Empty tests getting the default branch for a repo
func Test_GetDefaultBranch_Empty(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposByOwnerByRepo,
			&github.Repository{
				DefaultBranch: github.String(""),
			},
		))
	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	result, err := gh.GetDefaultBranch(testCommitModel.Owner, testCommitModel.Repo)
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
}

// Test_LogRateLimit tests that logging rate limit information is successful
func Test_LogRateLimit(t *testing.T) {
	assert.NotPanics(t, func() {
		logRateLimit(&github.Response{
			Rate: github.Rate{
				Limit:     100,
				Remaining: 100,
				Reset:     github.Timestamp{Time: time.Now()},
			},
		})
	})
}

// Test_ParseGithubWorkflowRuns tests parsing testWorkflowRun into testWorkflowRunModel
func Test_ParseGithubWorkflowRuns(t *testing.T) {
	parsed := parseGithubWorkflowRuns(testCommitModel.Owner, testCommitModel.Repo, testCommitModel.Branch, []*github.WorkflowRun{testWorkflowRun})
	assert.Equal(t, []*model.WorkflowRun{testWorkflowRunModel}, parsed)
}

// Test_GetWorkflowRuns tests getting all the test workflow runs
func Test_GetWorkflowRuns(t *testing.T) {

	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchPages(
			mock.GetReposActionsRunsByOwnerByRepo,
			&github.WorkflowRuns{
				TotalCount: github.Int(2),
				WorkflowRuns: []*github.WorkflowRun{
					testWorkflowRun, testWorkflowRun,
				},
			},
			&github.WorkflowRuns{
				TotalCount: github.Int(2),
				WorkflowRuns: []*github.WorkflowRun{
					testWorkflowRun, testWorkflowRun,
				},
			},
		),
	)

	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	results, err := gh.GetWorkflowRuns("owner1", "repo1", "branch1")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(results))
}

// Test_GetWorkflowRunAttempt tests getting a workflow run attempt
func Test_GetWorkflowRunAttempt(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			GetWorkflowRunAttempt,
			&testWorkflowRun,
		))
	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	result, err := gh.GetWorkflowRunAttempt(testCommitModel.Owner, testCommitModel.Repo, testCommitModel.Branch, 333, 2)
	assert.NoError(t, err)
	assert.Equal(t, testWorkflowRunModel, result)
}

// Test_GetContents tests getting contents of a file from a commit
func Test_GetContents(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposContentsByOwnerByRepoByPath,
			&github.RepositoryContent{
				Content: github.String("test"),
			},
		))
	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	result, err := gh.GetContents(testCommitModel.Owner, testCommitModel.Repo, "main.tf", "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "test", result)
}

// Test_adapter_GetPullRequests tests getting all pull requests in a specific state
func Test_adapter_GetPullRequests(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchPages(
			mock.GetReposPullsByOwnerByRepo,
			[]*github.PullRequest{testPullRequest, testPullRequest},
			[]*github.PullRequest{testPullRequest, testPullRequest},
		),
	)

	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	results, err := gh.GetPullRequests(owner, repo, "all")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(results))
}

// Test_parsePullRequests tests parsing testPullRequest into testPullRequestModel
func Test_parsePullRequests(t *testing.T) {
	parsed := parsePullRequests(testCommitModel.Owner, testCommitModel.Repo, []*github.PullRequest{testPullRequest})
	assert.Equal(t, []*model.PullRequest{testPullRequestModel}, parsed)
}

// Test_adapter_GetPullRequestCommits tests getting all pull request commits for a specific pull request
func Test_adapter_GetPullRequestCommits(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchPages(
			mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
			[]*github.RepositoryCommit{testCommitGitHub, testCommitGitHub},
			[]*github.RepositoryCommit{testCommitGitHub, testCommitGitHub},
		),
	)

	client := github.NewClient(mockedHTTPClient)
	gh := &adapter{client: client}
	results, err := gh.GetPullRequestCommits("owner1", "repo1", 10)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(results))
}
