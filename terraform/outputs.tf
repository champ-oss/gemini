output "region" {
  description = "AWS Region"
  value       = data.aws_region.this.name
}

output "db_name" {
  description = "Name of database"
  value       = aws_rds_cluster.this.database_name
}

output "db_arn" {
  description = "RDS cluster ARN"
  value       = aws_rds_cluster.this.arn
}

output "db_secrets_arn" {
  depends_on  = [aws_rds_cluster.this]
  description = "AWS Secrets ARN"
  value       = aws_secretsmanager_secret.this.arn
}

output "db_endpoint" {
  description = "DNS hostname of gemini database server"
  value       = aws_rds_cluster.this.endpoint
}

output "db_port" {
  description = "TCP port of gemini database server"
  value       = aws_rds_cluster.this.port
}

output "db_username" {
  description = "Username to connect to the gemini database server"
  depends_on  = [aws_rds_cluster.this]
  value       = var.database_username
}

output "db_password" {
  description = "Password to connect to the gemini database server"
  depends_on  = [aws_rds_cluster.this]
  sensitive   = true
  value       = random_password.database.result
}

output "grafana_dns" {
  description = "Grafana DNS hostname"
  depends_on  = [module.grafana]
  value       = local.grafana_dns
}

output "grafana_username" {
  description = "Grafana admin username"
  depends_on  = [module.grafana]
  value       = var.grafana_username
}

output "grafana_password" {
  description = "Grafana admin password"
  depends_on  = [module.grafana]
  sensitive   = true
  value       = random_password.grafana.result
}