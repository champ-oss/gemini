package model

import "gorm.io/gorm"

type Repo struct {
	gorm.Model
	Owner  string
	Name   string
	Branch string
}
