data "aws_kms_secrets" "this" {
  secret {
    name    = "terraform_api_key"
    payload = var.terraform_api_key
  }
}

provider "grafana" {
  url     = "https://${local.grafana_dns}"
  auth    = var.use_terraform_api_key ? data.aws_kms_secrets.this.plaintext["terraform_api_key"] : "${var.grafana_username}:${random_password.grafana.result}"
  retries = 100
}

# This will block the creation of resources until Grafana is fully up and running
resource "null_resource" "wait_for_grafana" {
  depends_on = [module.grafana]
  provisioner "local-exec" {
    command     = "curl --silent --fail --retry 180 --retry-delay 5 --retry-connrefused --insecure $URL/api/health"
    interpreter = ["/bin/sh", "-c"]
    environment = {
      URL = "https://${local.grafana_dns}"
    }
  }
}

resource "grafana_data_source" "this" {
  depends_on    = [null_resource.wait_for_grafana]
  type          = "mysql"
  name          = "gemini"
  is_default    = true
  url           = "${aws_rds_cluster.this.endpoint}:${aws_rds_cluster.this.port}"
  username      = var.database_username
  password      = random_password.database.result
  database_name = aws_rds_cluster.this.database_name
  secure_json_data {
    access_key = ""
    secret_key = ""
    password   = random_password.database.result # This is needed as setting the password above does not seem to work
  }
}

resource "grafana_dashboard" "status" {
  depends_on  = [null_resource.wait_for_grafana]
  config_json = file("${path.module}/gemini_status.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "deployment_frequency" {
  depends_on  = [null_resource.wait_for_grafana]
  config_json = file("${path.module}/deployment_frequency.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "change_failures" {
  depends_on  = [null_resource.wait_for_grafana]
  config_json = file("${path.module}/change_failures.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "lead_time_for_changes" {
  depends_on  = [null_resource.wait_for_grafana]
  config_json = file("${path.module}/lead_time_for_changes.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "time_to_restore" {
  depends_on  = [null_resource.wait_for_grafana]
  config_json = file("${path.module}/time_to_restore.json")
  message     = "Updated by Terraform"
  overwrite   = true
}