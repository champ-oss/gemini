resource "random_password" "database" {
  length  = 32
  special = false

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_db_subnet_group" "this" {
  name_prefix = "${var.git}-${random_string.identifier.result}-"
  subnet_ids  = var.private_subnet_ids
  tags        = var.tags

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "rds" {
  name_prefix = "${var.git}-${random_string.identifier.result}-"
  vpc_id      = var.vpc_id
  tags        = var.tags

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "rds_ingress_ecs" {
  description              = "ECS"
  type                     = "ingress"
  from_port                = aws_rds_cluster.this.port
  to_port                  = aws_rds_cluster.this.port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rds.id
  source_security_group_id = module.core.ecs_app_security_group
}

resource "aws_rds_cluster" "this" {
  cluster_identifier_prefix       = "${var.git}-${random_string.identifier.result}-"
  final_snapshot_identifier       = "${var.git}-${formatdate("YYYYMMDDhhmmss", timestamp())}"
  copy_tags_to_snapshot           = true
  engine                          = "aurora-mysql"
  engine_mode                     = "serverless"
  engine_version                  = "5.7.mysql_aurora.2.07.1" # cannot be modified after creation
  database_name                   = "grafana"
  master_username                 = var.database_username
  master_password                 = random_password.database.result
  backup_retention_period         = 5 # days
  snapshot_identifier             = var.snapshot_identifier
  vpc_security_group_ids          = [aws_security_group.rds.id]
  db_subnet_group_name            = aws_db_subnet_group.this.id
  db_cluster_parameter_group_name = var.db_cluster_parameter_group_name
  deletion_protection             = var.protect
  enable_http_endpoint            = true
  tags                            = var.tags

  scaling_configuration {
    auto_pause   = var.database_auto_pause
    min_capacity = 1
    max_capacity = var.database_max_capacity
  }

  lifecycle {
    create_before_destroy = true
    ignore_changes = [
      snapshot_identifier,
      final_snapshot_identifier,
      engine_version
    ]
  }
}
