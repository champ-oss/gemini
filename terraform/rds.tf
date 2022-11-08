module "aurora" {
  source                    = "github.com/champ-oss/terraform-aws-aurora.git?ref=45c1ed44506e9b4c2a9006804c2ab06bb5babe41"
  private_subnet_ids        = var.private_subnet_ids
  vpc_id                    = module.core.vpc_id
  source_security_group_id  = module.core.ecs_app_security_group
  cluster_identifier_prefix = var.git
  database_name             = "grafana"
  master_username           = var.database_username
  backup_retention_period   = 5 # days
  protect                   = var.protect
  tags                      = merge(local.tags, var.tags)
  skip_final_snapshot       = false
  git                       = var.git
  max_capacity              = var.database_max_capacity
}