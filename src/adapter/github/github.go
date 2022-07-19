package github

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/champ-oss/gemini/model"
	"github.com/google/go-github/v42/github"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	maxPerPage                 = 100 // 100 is max
	waitSecondsBetweenRequests = 2   // throttle requests to GitHub to avoid rate limiting
)

type adapter struct {
	client *github.Client
}

// NewAdapter initializes a new adapter
func NewAdapter(appID, installationID int64, pem []byte) (*adapter, error) {
	if len(pem) == 0 || appID == 0 || installationID == 0 {
		log.Warn("No GitHub credential details specified. Using unauthenticated client.")
		return &adapter{
			client: github.NewClient(nil),
		}, nil
	}

	transport, err := ghinstallation.New(http.DefaultTransport, appID, installationID, pem)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Transport: transport}
	client := github.NewClient(httpClient)

	return &adapter{
		client: client,
	}, nil
}

// GetCommits fetches all commits for the given repo and branch
func (a *adapter) GetCommits(owner string, repo string, branch string) ([]*model.Commit, error) {
	return a.GetCommitsForFile(owner, repo, branch, "")
}

func (a *adapter) GetCommitsForFile(owner string, repo string, branch string, filePath string) ([]*model.Commit, error) {
	options := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: maxPerPage},
		SHA:         branch,
		Path:        filePath,
	}

	var allCommits []*github.RepositoryCommit

	for {
		commits, resp, err := a.client.Repositories.ListCommits(context.Background(), owner, repo, options)
		time.Sleep(time.Second * waitSecondsBetweenRequests)
		if err != nil {
			return nil, err
		}
		logRateLimit(resp)

		allCommits = append(allCommits, commits...)
		log.Debug(fmt.Sprintf("Retrieved %d commits from GitHub", len(allCommits)))
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return parseGithubCommits(owner, repo, branch, allCommits), nil
}

// parseGithubCommits converts GitHub commits to our standard commit model
func parseGithubCommits(owner string, repo string, branch string, commits []*github.RepositoryCommit) []*model.Commit {
	var parsed []*model.Commit
	for _, commit := range commits {
		parsed = append(parsed, &model.Commit{
			Owner:          owner,
			Repo:           repo,
			Branch:         branch,
			Message:        commit.Commit.GetMessage(),
			CommitterName:  commit.Commit.Committer.GetName(),
			CommitterDate:  commit.Commit.Committer.Date.Unix(),
			CommitterEmail: commit.Commit.Committer.GetEmail(),
			AuthorName:     commit.Commit.Author.GetName(),
			AuthorDate:     commit.Commit.Author.Date.Unix(),
			AuthorEmail:    commit.Commit.Author.GetEmail(),
			Url:            commit.GetURL(),
			Sha:            commit.GetSHA(),
		})
	}
	return parsed
}

// GetDefaultBranch queries a repository and returns the name of the default branch (ex: main)
func (a *adapter) GetDefaultBranch(owner string, repo string) (string, error) {
	output, resp, err := a.client.Repositories.Get(context.Background(), owner, repo)
	time.Sleep(time.Second * waitSecondsBetweenRequests)
	if err != nil {
		return "", err
	}
	logRateLimit(resp)
	if *output.DefaultBranch == "" {
		return "", fmt.Errorf("no default branch set")
	}
	return *output.DefaultBranch, nil
}

// logRateLimit logs the GitHub rate limit information
func logRateLimit(resp *github.Response) {
	log.Debugf("GitHub Rate Limit: %d Remaining: %d Reset: %s", resp.Rate.Limit, resp.Rate.Remaining, resp.Rate.Reset)
}

// GetWorkflowRuns fetches all workflow runs for the given repo and branch
func (a *adapter) GetWorkflowRuns(owner string, repo string, branch string) ([]*model.WorkflowRun, error) {
	options := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: maxPerPage},
		Branch:      branch,
		Status:      "completed",
	}

	var allRuns []*github.WorkflowRun

	for {
		runs, resp, err := a.client.Actions.ListRepositoryWorkflowRuns(context.Background(), owner, repo, options)
		time.Sleep(time.Second * waitSecondsBetweenRequests)
		if err != nil {
			return nil, err
		}
		logRateLimit(resp)

		allRuns = append(allRuns, runs.WorkflowRuns...)
		log.Debug(fmt.Sprintf("Retrieved %d workflow runs from GitHub", len(allRuns)))

		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return parseGithubWorkflowRuns(owner, repo, branch, allRuns), nil
}

// GetWorkflowRunAttempt gets a specific workflow run by run ID and attempt number
func (a *adapter) GetWorkflowRunAttempt(owner string, repo string, branch string, runID int64, attemptNumber int) (*model.WorkflowRun, error) {
	run, resp, err := a.client.Actions.GetWorkflowRunAttempt(context.Background(), owner, repo, runID, attemptNumber, nil)
	time.Sleep(time.Second * waitSecondsBetweenRequests)
	if err != nil {
		return nil, err
	}
	logRateLimit(resp)
	parsed := parseGithubWorkflowRuns(owner, repo, branch, []*github.WorkflowRun{run})
	return parsed[0], nil
}

// parseGithubWorkflowRuns converts GitHub workflow runs to our standard WorkflowRun model
func parseGithubWorkflowRuns(owner string, repo string, branch string, runs []*github.WorkflowRun) []*model.WorkflowRun {
	var parsed []*model.WorkflowRun
	for _, run := range runs {
		parsed = append(parsed, &model.WorkflowRun{
			Owner:        owner,
			Repo:         repo,
			Branch:       branch,
			NodeID:       run.GetNodeID(),
			Name:         run.GetName(),
			Sha:          run.GetHeadSHA(),
			Conclusion:   run.GetConclusion(),
			RunCreatedAt: run.CreatedAt.Unix(),
			RunUpdatedAt: run.UpdatedAt.Unix(),
			RunAttempt:   run.GetRunAttempt(),
			RunID:        run.GetID(),
		})
	}
	return parsed
}

// GetContents retrieves a file from Github from a specific commit sha
func (a *adapter) GetContents(owner string, repo string, filePath string, ref string) (string, error) {
	contents, _, resp, err := a.client.Repositories.GetContents(context.Background(), owner, repo, filePath, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	time.Sleep(time.Second * waitSecondsBetweenRequests)
	logRateLimit(resp)
	if err != nil {
		return "", err
	}
	return contents.GetContent()
}

// GetPullRequests fetches all pull requests for the given repo
func (a *adapter) GetPullRequests(owner string, repo string, state string) ([]*model.PullRequest, error) {
	options := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: maxPerPage},
		State:       state,
	}

	var allPullRequests []*github.PullRequest

	for {
		pullRequests, resp, err := a.client.PullRequests.List(context.Background(), owner, repo, options)
		time.Sleep(time.Second * waitSecondsBetweenRequests)
		if err != nil {
			return nil, err
		}
		logRateLimit(resp)

		allPullRequests = append(allPullRequests, pullRequests...)
		log.Debug(fmt.Sprintf("Retrieved %d pull requests from GitHub", len(allPullRequests)))

		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return parsePullRequests(owner, repo, allPullRequests), nil
}

// parsePullRequests converts GitHub pull requests to our standard PullRequest model
func parsePullRequests(owner string, repo string, pullRequests []*github.PullRequest) []*model.PullRequest {
	var parsed []*model.PullRequest
	for _, pullRequest := range pullRequests {
		parsed = append(parsed, &model.PullRequest{
			Owner:                owner,
			Repo:                 repo,
			MergeCommitSHA:       pullRequest.GetMergeCommitSHA(),
			Number:               pullRequest.GetNumber(),
			State:                pullRequest.GetState(),
			Title:                pullRequest.GetTitle(),
			PullRequestCreatedAt: pullRequest.GetCreatedAt().Unix(),
			PullRequestUpdatedAt: pullRequest.GetUpdatedAt().Unix(),
			PullRequestClosedAt:  pullRequest.GetClosedAt().Unix(),
			PullRequestMergedAt:  pullRequest.GetMergedAt().Unix(),
			Draft:                pullRequest.GetDraft(),
			Merged:               pullRequest.GetMerged(),
			Commits:              pullRequest.GetCommits(),
		})
	}
	return parsed
}

// GetPullRequestCommits fetches all pull requests commits for the given repo
func (a *adapter) GetPullRequestCommits(owner string, repo string, number int) ([]*model.Commit, error) {
	options := &github.ListOptions{PerPage: maxPerPage}

	var allCommits []*github.RepositoryCommit

	for {
		commits, resp, err := a.client.PullRequests.ListCommits(context.Background(), owner, repo, number, options)
		time.Sleep(time.Second * waitSecondsBetweenRequests)
		if err != nil {
			return nil, err
		}
		logRateLimit(resp)

		allCommits = append(allCommits, commits...)
		log.Debug(fmt.Sprintf("Retrieved %d pull request commits from GitHub", len(allCommits)))
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return parseGithubCommits(owner, repo, "", allCommits), nil
}
