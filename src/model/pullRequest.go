package model

import "gorm.io/gorm"

type PullRequest struct {
	gorm.Model
	Owner                string
	Repo                 string `gorm:"index"`
	Number               int
	MergeCommitSHA       string `gorm:"unique"`
	State                string
	Title                string
	PullRequestCreatedAt int64
	PullRequestUpdatedAt int64
	PullRequestClosedAt  int64
	PullRequestMergedAt  int64
	Draft                bool
	Merged               bool
	Commits              int
}
