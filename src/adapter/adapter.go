package adapter

import "github.com/champ-oss/gemini/model"

// GitClient represents a generic interface for accessing information for a Git repository
type GitClient interface {
	GetCommits(owner string, repo string, branch string) ([]*model.Commit, error)
	GetCommitsForFile(owner string, repo string, branch string, filePath string) ([]*model.Commit, error)
	GetDefaultBranch(owner string, repo string) (string, error)
	GetWorkflowRuns(owner string, repo string, branch string) ([]*model.WorkflowRun, error)
	GetWorkflowRunAttempt(owner string, repo string, branch string, runID int64, attemptNumber int) (*model.WorkflowRun, error)
	GetContents(owner string, repo string, filePath string, ref string) (string, error)
	GetPullRequests(owner string, repo string, state string) ([]*model.PullRequest, error)
	GetPullRequestCommits(owner string, repo string, number int) ([]*model.Commit, error)
}
