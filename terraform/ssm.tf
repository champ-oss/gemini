resource "aws_ssm_parameter" "database" {
  name        = "${var.git}-${random_string.identifier.result}-database-password"
  description = "Database password"
  type        = "SecureString"
  value       = random_password.database.result
  tags        = var.tags
}

resource "aws_ssm_parameter" "grafana" {
  name        = "${var.git}-${random_string.identifier.result}-grafana-password"
  description = "Grafana password"
  type        = "SecureString"
  value       = random_password.grafana.result
  tags        = var.tags
}