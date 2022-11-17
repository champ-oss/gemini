package repository

import (
	"fmt"
	"github.com/champ-oss/gemini/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository interface {
	AddCommits(commits []*model.Commit) (inserted int64, err error)
	AddWorkflowRuns(runs []*model.WorkflowRun) (inserted int64, err error)
	AddTerraformRefs(runs []*model.TerraformRef) (inserted int64, err error)
	GetWorkflowRunsByName(owner string, repo string, branch string, name string) []*model.WorkflowRun
	AddPullRequestCommits(pullRequests []*model.PullRequestCommit) (inserted int64, err error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository initializes a new repository
func NewRepository(username, password, hostname, port, database string, dropTables bool, db *gorm.DB) (*repository, error) {
	if db == nil {
		log.WithFields(log.Fields{"username": username, "hostname": hostname, "port": port, "database": database}).Info("Connecting to database")
		dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, database)
		var err error
		db, err = gorm.Open(mysql.Open(dataSource), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}

	repo := &repository{
		db,
	}

	if dropTables {
		log.Warn("Dropping existing database tables")
		dropDatabaseTables(repo)
	}

	err := initializeDatabase(repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func initializeDatabase(repo *repository) error {
	err := repo.db.AutoMigrate(
		&model.Commit{},
		&model.WorkflowRun{},
		&model.TerraformRef{},
		&model.PullRequestCommit{},
	)
	if err != nil {
		return err
	}
	return nil
}

func dropDatabaseTables(repo *repository) {
	_ = repo.db.Migrator().DropTable(&model.Commit{})
	_ = repo.db.Migrator().DropTable(&model.WorkflowRun{})
	_ = repo.db.Migrator().DropTable(&model.TerraformRef{})
	_ = repo.db.Migrator().DropTable(&model.PullRequestCommit{})
}

func (r *repository) AddCommits(commits []*model.Commit) (inserted int64, err error) {
	if len(commits) == 0 {
		return 0, nil
	}
	result := r.db.Model(&model.Commit{}).Save(commits)
	return result.RowsAffected, result.Error
}

func (r *repository) AddWorkflowRuns(runs []*model.WorkflowRun) (inserted int64, err error) {
	if len(runs) == 0 {
		return 0, nil
	}
	result := r.db.Model(&model.WorkflowRun{}).Save(runs)
	return result.RowsAffected, result.Error
}

func (r *repository) AddTerraformRefs(runs []*model.TerraformRef) (inserted int64, err error) {
	if len(runs) == 0 {
		return 0, nil
	}
	result := r.db.Model(&model.TerraformRef{}).Save(runs)
	return result.RowsAffected, result.Error
}

func (r *repository) GetWorkflowRunsByName(owner string, repo string, branch string, name string) []*model.WorkflowRun {
	var runs []*model.WorkflowRun
	r.db.Model(&model.WorkflowRun{}).
		Where(&model.WorkflowRun{Owner: owner, Repo: repo, Branch: branch}).
		Where("name LIKE ?", "%"+name+"%").
		Find(&runs)
	return runs
}

func (r *repository) AddPullRequestCommits(pullRequests []*model.PullRequestCommit) (inserted int64, err error) {
	if len(pullRequests) == 0 {
		return 0, nil
	}
	result := r.db.Model(&model.PullRequestCommit{}).Save(pullRequests)
	return result.RowsAffected, result.Error
}
