data "aws_kms_secrets" "this" {
  count = var.use_terraform_api_key ? 1 : 0
  secret {
    name    = "terraform_api_key"
    payload = var.terraform_api_key
  }
}

provider "grafana" {
  url     = "https://${local.grafana_dns}"
  auth    = var.use_terraform_api_key ? data.aws_kms_secrets.this[0].plaintext["terraform_api_key"] : "${var.grafana_username}:${random_password.grafana.result}"
  retries = 300
}

resource "grafana_data_source" "this" {
  depends_on    = [module.grafana, module.aurora]
  type          = "mysql"
  name          = "gemini"
  is_default    = true
  url           = "${module.aurora.endpoint}:${module.aurora.port}"
  username      = var.database_username
  password      = module.aurora.master_password
  database_name = module.aurora.database_name
  secure_json_data {
    access_key = ""
    secret_key = ""
    password   = module.aurora.master_password # This is needed as setting the password above does not seem to work
  }
}

resource "grafana_dashboard" "status" {
  config_json = file("${path.module}/gemini_status.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "deployment_frequency" {
  config_json = file("${path.module}/deployment_frequency.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "change_failures" {
  config_json = file("${path.module}/change_failures.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "lead_time_for_changes" {
  config_json = file("${path.module}/lead_time_for_changes.json")
  message     = "Updated by Terraform"
  overwrite   = true
}

resource "grafana_dashboard" "time_to_restore" {
  config_json = file("${path.module}/time_to_restore.json")
  message     = "Updated by Terraform"
  overwrite   = true
}
