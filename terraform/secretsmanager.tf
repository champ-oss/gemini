resource "aws_secretsmanager_secret" "this" {
  name                    = "rds-db-credentials/${aws_rds_cluster.this.cluster_resource_id}/${var.git}-${random_string.identifier.result}"
  description             = "RDS credentials for use in query editor"
  recovery_window_in_days = 0
  tags                    = var.tags
}

locals {
  secret = {
    dbInstanceIdentifier = aws_rds_cluster.this.cluster_identifier
    engine               = aws_rds_cluster.this.engine
    dbname               = aws_rds_cluster.this.database_name
    host                 = aws_rds_cluster.this.endpoint
    port                 = aws_rds_cluster.this.port
    resourceId           = aws_rds_cluster.this.cluster_resource_id
    username             = aws_rds_cluster.this.master_username
    password             = random_password.database.result
  }
}

resource "aws_secretsmanager_secret_version" "this" {
  secret_id     = aws_secretsmanager_secret.this.id
  secret_string = jsonencode(local.secret)
}