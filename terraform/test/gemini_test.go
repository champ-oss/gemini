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
	defer destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	grafanaDns := terraform.Output(t, terraformOptions, "grafana_dns")
	grafanaUsername := terraform.Output(t, terraformOptions, "grafana_username")
	grafanaPassword := terraform.Output(t, terraformOptions, "grafana_password")

	// Verify that all database tables are populated with data
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "commits", 11))
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "workflow_runs", 74))
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "terraform_refs", 3))
	assert.NoError(t, checkExpectedGrafanaTableCount(grafanaDns, grafanaUsername, grafanaPassword, "pull_request_commits", 7))
}

func destroy(t *testing.T, options *terraform.Options) {
	t.Log("removing grafana dashboard resources from state")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_data_source.this")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.status")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.deployment_frequency")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.change_failures")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.lead_time_for_changes")
	_, _ = terraform.RunTerraformCommandE(t, options, "state", "rm", "module.this.grafana_dashboard.time_to_restore")
}
