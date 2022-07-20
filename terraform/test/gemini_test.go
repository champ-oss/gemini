package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestGemini tests the application in an ephemeral environment
func TestGemini(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir:  "../examples/complete",
		BackendConfig: map[string]interface{}{},
		Vars: map[string]interface{}{
			"docker_tag":             os.Getenv("GITHUB_SHA"),
			"github_app_id":          os.Getenv("APP_ID"),
			"github_installation_id": os.Getenv("INSTALLATION_ID"),
			"github_pem":             os.Getenv("PEM"),
		},
	}
	defer destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	commitsTest(t, terraformOptions)
	grafanaTest(t, terraformOptions)
	actionsTest(t, terraformOptions)
	terraformrefsTest(t, terraformOptions)
	pullrequestsTest(t, terraformOptions)
}

// Validate that the table for commits is populated successfully
func commitsTest(t *testing.T, options *terraform.Options) {
	dbName := terraform.Output(t, options, "db_name")
	dbArn := terraform.Output(t, options, "db_arn")
	dbSecretsArn := terraform.Output(t, options, "db_secrets_arn")

	table := "commits"
	owner := "champtitles"
	repo := "tflint-ruleset-champtitles"
	expectedRows := 15

	awsSess := getAWSSession()
	defer dropTable(awsSess, dbName, dbArn, dbSecretsArn, table)

	t.Logf("Checking that the %s table is successfully fully populated on first run", table)
	assert.Nil(t, countRecords(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, expectedRows))

	t.Log("Testing removing some of the most recent commits from the DB so the process will sync them again")
	assert.Nil(t, deleteRecentCommits(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, 5))

	t.Log("Checking that commits are fully synced again after the process re-runs")
	assert.Nil(t, countRecords(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, expectedRows))

	t.Log("Checking that no records have empty fields")
	assert.Nil(t, checkForEmptyFields(awsSess, dbName, dbArn, dbSecretsArn, table))
}

// Validate that grafana is working correctly
func grafanaTest(t *testing.T, options *terraform.Options) {
	grafanaDns := terraform.Output(t, options, "grafana_dns")
	username := terraform.Output(t, options, "grafana_username")
	password := terraform.Output(t, options, "grafana_password")
	table := "workflow_runs"

	t.Logf("Running query to check Grafana data source connection. Host: %s", grafanaDns)

	queryUrl := fmt.Sprintf("https://%s:%s@%s/api/tsdb/query", username, password, grafanaDns)

	queryRequest := &GrafanaQueryRequest{
		From: "now-10y",
		To:   "now",
		Queries: []*GrafanaQuery{
			{
				IntervalMs:    int64(86400000),
				MaxDataPoints: 1000,
				DatasourceId:  1,
				RawSql:        fmt.Sprintf("SELECT count(*) FROM %s", table),
				Format:        "table",
			},
		},
	}

	queryRequestJson, err := json.Marshal(queryRequest)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", queryUrl, bytes.NewBuffer(queryRequestJson))
	assert.Nil(t, err)

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Content-Type", `application/json`)

	c := http.Client{Timeout: time.Duration(10) * time.Second}
	resp, err := c.Do(req)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	t.Logf("query returned: %s message: %s", resp.Status, body)
	assert.Nil(t, err)
	assert.Equal(t, "200 OK", resp.Status)
}

// Validate that the table for workflow runs is populated successfully
func actionsTest(t *testing.T, options *terraform.Options) {
	dbName := terraform.Output(t, options, "db_name")
	dbArn := terraform.Output(t, options, "db_arn")
	dbSecretsArn := terraform.Output(t, options, "db_secrets_arn")

	table := "workflow_runs"
	owner := "champtitles"
	repo := "tflint-ruleset-champtitles"
	expectedRows := 218

	awsSess := getAWSSession()
	defer dropTable(awsSess, dbName, dbArn, dbSecretsArn, table)

	t.Logf("Checking that the %s table is successfully fully populated on first run", table)
	assert.Nil(t, countRecords(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, expectedRows))

	t.Log("Checking that no records have empty fields")
	assert.Nil(t, checkForEmptyFields(awsSess, dbName, dbArn, dbSecretsArn, table))

	t.Log("Checking workflow reruns in the database")
	assert.Nil(t, validateReruns(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, 2))
}

// Validate that the table for terraform module references is populated successfully
func terraformrefsTest(t *testing.T, options *terraform.Options) {
	dbName := terraform.Output(t, options, "db_name")
	dbArn := terraform.Output(t, options, "db_arn")
	dbSecretsArn := terraform.Output(t, options, "db_secrets_arn")

	table := "terraform_refs"
	owner := "champtitles"
	repo := "tflint-ruleset-champtitles"
	expectedRows := 21

	awsSess := getAWSSession()
	defer dropTable(awsSess, dbName, dbArn, dbSecretsArn, table)

	t.Logf("Checking that the %s table is successfully fully populated on first run", table)
	assert.Nil(t, countRecords(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, expectedRows))

	t.Log("Checking that no records have empty fields")
	assert.Nil(t, checkForEmptyFields(awsSess, dbName, dbArn, dbSecretsArn, table))
}

// Validate that the table for pull request commits is populated successfully
func pullrequestsTest(t *testing.T, options *terraform.Options) {
	dbName := terraform.Output(t, options, "db_name")
	dbArn := terraform.Output(t, options, "db_arn")
	dbSecretsArn := terraform.Output(t, options, "db_secrets_arn")

	table := "pull_request_commits"
	owner := "champtitles"
	repo := "tflint-ruleset-champtitles"
	expectedRows := 26

	awsSess := getAWSSession()
	defer dropTable(awsSess, dbName, dbArn, dbSecretsArn, table)

	t.Logf("Checking that the %s table is successfully fully populated on first run", table)
	assert.Nil(t, countRecords(awsSess, dbName, dbArn, dbSecretsArn, table, owner, repo, expectedRows))
}

func destroy(t *testing.T, options *terraform.Options) {

	t.Log("removing grafana dashboard resources from state")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_data_source.this")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.status")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.deployment_frequency")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.change_failures")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.lead_time_for_changes")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.time_to_restore")

	t.Log("Running Terraform Destroy")
	terraform.Destroy(t, options)
}
