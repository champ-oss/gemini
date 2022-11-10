output "region" {
  description = "AWS Region"
  value       = data.aws_region.this.name
}

output "db_name" {
  description = "Name of database"
  value       = module.aurora.database_name
}

output "db_arn" {
  description = "RDS cluster ARN"
  value       = module.aurora.arn
}

output "db_endpoint" {
  description = "DNS hostname of gemini database server"
  value       = module.aurora.endpoint
}

output "db_port" {
  description = "TCP port of gemini database server"
  value       = module.aurora.port
}

output "db_username" {
  description = "Username to connect to the gemini database server"
  depends_on  = [module.aurora]
  value       = var.database_username
}

output "db_password" {
  description = "Password to connect to the gemini database server"
  depends_on  = [module.aurora]
  sensitive   = true
  value       = module.aurora.master_password
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

output "hash" {
  value = module.hash.hash
}