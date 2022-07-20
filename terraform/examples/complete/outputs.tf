output "region" {
  description = "AWS Region"
  value       = module.this.region
}

output "db_name" {
  description = "Name of database"
  value       = module.this.db_name
}

output "db_arn" {
  description = "RDS cluster ARN"
  value       = module.this.db_arn
}

output "db_secrets_arn" {
  description = "AWS Secrets ARN"
  value       = module.this.db_secrets_arn
}

output "grafana_dns" {
  description = "Grafana DNS hostname"
  value       = module.this.grafana_dns
}

output "grafana_username" {
  description = "Grafana admin username"
  value       = module.this.grafana_username
}

output "grafana_password" {
  description = "Grafana admin password"
  sensitive   = true
  value       = module.this.grafana_password
}