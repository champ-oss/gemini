package model

import "gorm.io/gorm"

type WorkflowRun struct {
	gorm.Model
	Owner        string
	Repo         string `gorm:"index"`
	Name         string `gorm:"index"`
	Branch       string
	Sha          string
	Conclusion   string `gorm:"index"`
	RunCreatedAt int64
	RunUpdatedAt int64
	RunAttempt   int    `gorm:"uniqueIndex:idx_nodeid_runattempt"`
	NodeID       string `gorm:"uniqueIndex:idx_nodeid_runattempt;size:256"`
	RunID        int64
}
