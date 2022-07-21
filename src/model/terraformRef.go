package model

import "gorm.io/gorm"

type TerraformRef struct {
	gorm.Model
	Owner        string
	Repo         string `gorm:"index"`
	Branch       string
	Sha          string `gorm:"uniqueIndex:idx_sha_module_name;size:256"`
	FileName     string
	RunUpdatedAt int64
	ModuleName   string `gorm:"uniqueIndex:idx_sha_module_name;size:256"`
	SourceOwner  string
	SourceRepo   string
	SourceRef    string
}
