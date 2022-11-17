package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champ-oss/gemini/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var TestCommit1 = &model.Commit{
	Owner:          "owner1",
	Repo:           "repo1",
	Branch:         "branch1",
	Message:        "commit msg",
	CommitterName:  "test committer",
	CommitterDate:  1642708819,
	CommitterEmail: "testcommitter@test.com",
	AuthorName:     "test author",
	AuthorDate:     1642708821,
	AuthorEmail:    "testauthor@test.com",
	Url:            "http://localhost",
	Sha:            "702fe8ee76f422edd4bc257a0a2171af26563063",
}

var TestCommit2 = &model.Commit{
	Owner:  "owner2",
	Repo:   "repo2",
	Branch: "branch2",
}

var TestWorkflowRun1 = &model.WorkflowRun{
	Owner:  "owner2",
	Repo:   "repo2",
	Branch: "branch2",
}

var TestTerraformRef = &model.TerraformRef{
	Repo:         "repo3",
	Sha:          "702fe8ee76f422edd4bc257a0a2171af26563063",
	FileName:     "main.tf",
	RunUpdatedAt: 1642708821,
	ModuleName:   "this",
	SourceOwner:  "owner1",
	SourceRepo:   "repo1",
	SourceRef:    "702fe8ee721398hr19j16563063",
}

var testPullRequestCommitModel = &model.PullRequestCommit{
	Owner:          "owner2",
	Repo:           "repo2",
	MergeCommitSHA: "702fe8ee76f422edd4bc257a0a2171af26563063",
	Number:         10,
	Sha:            "702fe8ee76f422edd4bc257a0a2171af26563063",
}

func getMockRepo() (repository, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{SkipDefaultTransaction: true})

	return repository{gormDB}, mock
}

func Test_NewRepository(t *testing.T) {
	_, err := NewRepository("", "", "", "", "", false)
	assert.Error(t, err)
}

func Test_InitializeDatabase(t *testing.T) {
	repo, mock := getMockRepo()
	mock.ExpectQuery("SELECT SCHEMA_NAME from Information_schema.SCHEMATA").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectExec("CREATE TABLE `commits`").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery("SELECT SCHEMA_NAME from Information_schema.SCHEMATA").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectExec("CREATE TABLE `workflow_runs`").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery("SELECT SCHEMA_NAME from Information_schema.SCHEMATA").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectExec("CREATE TABLE `terraform_refs`").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery("SELECT SCHEMA_NAME from Information_schema.SCHEMATA").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectExec("CREATE TABLE `pull_request_commits`").WillReturnResult(sqlmock.NewResult(0, 0))
	err := initializeDatabase(&repo)
	assert.Nil(t, err)
}

func Test_InitializeDatabase_Error(t *testing.T) {
	repo, mock := getMockRepo()
	mock.ExpectQuery("SELECT SCHEMA_NAME from Information_schema.SCHEMATA").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectExec("CREATE TABLE `commits`").WillReturnError(fmt.Errorf("test error"))
	assert.Error(t, initializeDatabase(&repo))
}

func Test_DropDatabaseTables(t *testing.T) {
	repo, mock := getMockRepo()
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 0;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS `commits` CASCADE").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 1;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 0;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS `workflow_runs` CASCADE").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 1;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 0;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS `terraform_refs` CASCADE").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 1;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 0;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS `pull_request_commits` CASCADE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 1;").WillReturnResult(sqlmock.NewResult(0, 0))
	dropDatabaseTables(&repo)
}

func Test_AddCommits(t *testing.T) {
	repo, mock := getMockRepo()

	query := "INSERT INTO `commits` .*"

	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 2))

	inserted, err := repo.AddCommits([]*model.Commit{TestCommit1, TestCommit2})
	assert.Nil(t, err)
	assert.Equal(t, int64(2), inserted)
}

func Test_AddCommits_Empty(t *testing.T) {
	repo, _ := getMockRepo()
	inserted, err := repo.AddCommits([]*model.Commit{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), inserted)
}

func Test_AddWorkflowRuns(t *testing.T) {
	repo, mock := getMockRepo()

	query := "INSERT INTO `workflow_runs` .*"

	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))

	inserted, err := repo.AddWorkflowRuns([]*model.WorkflowRun{TestWorkflowRun1})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), inserted)
}

func Test_AddWorkflowRuns_Empty(t *testing.T) {
	repo, _ := getMockRepo()
	inserted, err := repo.AddWorkflowRuns([]*model.WorkflowRun{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), inserted)
}

func Test_AddTerraformRefs(t *testing.T) {
	repo, mock := getMockRepo()

	query := "INSERT INTO `terraform_refs` .*"

	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))

	inserted, err := repo.AddTerraformRefs([]*model.TerraformRef{TestTerraformRef})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), inserted)
}

func Test_AddTerraformRefs_Empty(t *testing.T) {
	repo, _ := getMockRepo()
	inserted, err := repo.AddTerraformRefs([]*model.TerraformRef{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), inserted)
}

func Test_GetWorkflowRunsByName(t *testing.T) {
	repo, mock := getMockRepo()

	query := "SELECT .* FROM `workflow_runs` .*"

	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"owner", "repo", "branch"}).AddRow("", "", ""))

	results := repo.GetWorkflowRunsByName(TestWorkflowRun1.Owner, TestWorkflowRun1.Repo, TestWorkflowRun1.Branch, "apply")
	assert.Len(t, results, 1)
}

func Test_repository_AddPullRequestCommits(t *testing.T) {
	repo, mock := getMockRepo()

	query := "INSERT INTO `pull_request_commits` .*"

	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))

	inserted, err := repo.AddPullRequestCommits([]*model.PullRequestCommit{testPullRequestCommitModel})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), inserted)
}

func Test_AddPullRequestCommits_Empty(t *testing.T) {
	repo, _ := getMockRepo()
	inserted, err := repo.AddPullRequestCommits([]*model.PullRequestCommit{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), inserted)
}
