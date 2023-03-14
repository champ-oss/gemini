module "aurora" {
  source                    = "github.com/champ-oss/terraform-aws-aurora.git?ref=v1.0.35-db258d6"
  private_subnet_ids        = var.private_subnet_ids
  vpc_id                    = var.vpc_id
  source_security_group_id  = module.core.ecs_app_security_group
  cluster_identifier_prefix = var.git
  cluster_instance_count    = var.cluster_instance_count
  database_name             = "grafana"
  master_username           = var.database_username
  backup_retention_period   = 5 # days
  protect                   = var.protect
  tags                      = merge(local.tags, var.tags)
  skip_final_snapshot       = false
  git                       = var.git
  max_capacity              = var.database_max_capacity
  metric_alarms_enabled     = var.metric_alarms_enabled
}