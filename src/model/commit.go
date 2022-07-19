package model

import "gorm.io/gorm"

type Commit struct {
	gorm.Model
	Owner          string
	Repo           string `gorm:"index"`
	Branch         string
	Message        string
	CommitterName  string
	CommitterDate  int64
	CommitterEmail string
	AuthorName     string
	AuthorDate     int64
	AuthorEmail    string
	Url            string
	Sha            string `gorm:"unique"`
}
