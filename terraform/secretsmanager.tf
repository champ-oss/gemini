resource "aws_secretsmanager_secret" "this" {
  name                    = "rds-db-credentials/${module.aurora.cluster_resource_id}/${var.git}-${random_string.identifier.result}"
  description             = "RDS credentials for use in query editor"
  recovery_window_in_days = 0
  tags                    = var.tags
}

locals {
  secret = {
    dbInstanceIdentifier = module.aurora.cluster_identifier
    engine               = "aurora-mysql"
    dbname               = module.aurora.database_name
    host                 = module.aurora.endpoint
    port                 = module.aurora.port
    resourceId           = module.aurora.cluster_resource_id
    username             = module.aurora.master_username
    password             = module.aurora.master_password
  }
}

resource "aws_secretsmanager_secret_version" "this" {
  secret_id     = aws_secretsmanager_secret.this.id
  secret_string = jsonencode(local.secret)
}