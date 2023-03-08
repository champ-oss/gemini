package test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestGemini tests the application in an ephemeral environment
func TestGemini(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/complete",
		BackendConfig: map[string]interface{}{
			"bucket": os.Getenv("TF_STATE_BUCKET"),
			"key":    os.Getenv("TF_VAR_git"),
		},
		Vars: map[string]interface{}{
			"github_app_id":          os.Getenv("APP_ID"),
			"github_installation_id": os.Getenv("INSTALLATION_ID"),
			"github_pem":             os.Getenv("PEM"),
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	grafanaDns := terraform.Output(t, terraformOptions, "grafana_dns")
	grafanaUsername := terraform.Output(t, terraformOptions, "grafana_username")
	grafanaPassword := terraform.Output(t, terraformOptions, "grafana_password")
	grafanaDataSourceId := terraform.Output(t, terraformOptions, "grafana_data_source_id")

	// Verify that all database tables are populated with data
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "commits", grafanaDataSourceId, 11))
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "workflow_runs", grafanaDataSourceId, 74))
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "terraform_refs", grafanaDataSourceId, 3))
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "pull_request_commits", grafanaDataSourceId, 7))
}
