package model

import "gorm.io/gorm"

type PullRequestCommit struct {
	gorm.Model
	Owner          string
	Repo           string `gorm:"index"`
	Number         int
	MergeCommitSHA string `gorm:"uniqueIndex:idx_sha_merge_commit_sha;size:256"`
	Sha            string `gorm:"uniqueIndex:idx_sha_merge_commit_sha;size:256"`
	CommitterDate  int64
}
