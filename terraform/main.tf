data "aws_region" "this" {}

locals {
  image_app = "champtitles/gemini:${module.hash.hash}"

  config_app = {
    DEBUG                  = var.debug ? "true" : "false"
    MINUTES_BETWEEN_CHECKS = var.minutes_between_checks
    DB_HOST                = module.aurora.endpoint
    DB_PORT                = tostring(module.aurora.port)
    DB_NAME                = module.aurora.database_name
    DB_USERNAME            = var.database_username
    REPOS                  = join(",", var.repos)
    DROP_TABLES            = var.drop_tables ? "true" : "false"
  }

  secrets = {
    DB_PASSWORD                = aws_ssm_parameter.database.name
    GF_DATABASE_PASSWORD       = aws_ssm_parameter.database.name
    GF_SECURITY_ADMIN_PASSWORD = aws_ssm_parameter.grafana.name
  }

  kms_secrets_app = {
    GITHUB_APP_ID          = var.github_app_id
    GITHUB_INSTALLATION_ID = var.github_installation_id
    GITHUB_PEM             = var.github_pem
  }

  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

module "hash" {
  source = "github.com/champ-oss/terraform-git-hash.git?ref=v1.0.13-ffd1b7d"
  path   = "${path.module}/.."
}

resource "random_string" "identifier" {
  length  = 5
  special = false
  upper   = false
  lower   = true
  numeric = true
}

module "core" {
  source                      = "github.com/champ-oss/terraform-aws-core.git?ref=v1.0.114-758b2d1"
  git                         = "${var.git}-${random_string.identifier.result}"
  name                        = "${var.git}-${random_string.identifier.result}"
  vpc_id                      = var.vpc_id
  public_subnet_ids           = var.public_subnet_ids
  private_subnet_ids          = var.private_subnet_ids
  certificate_arn             = var.certificate_arn
  additional_certificate_arns = var.additional_certificate_arns
  protect                     = false
  log_retention               = "731"
  tags                        = merge(local.tags, var.tags)
}

module "app" {
  source                = "github.com/champ-oss/terraform-aws-app?ref=v1.0.195-c710781"
  git                   = "${var.git}-${random_string.identifier.result}"
  vpc_id                = var.vpc_id
  subnets               = var.private_subnet_ids
  cluster               = module.core.ecs_cluster_name
  security_groups       = [module.core.ecs_app_security_group]
  execution_role_arn    = module.core.execution_ecs_role_arn
  wait_for_steady_state = var.wait_for_steady_state
  enable_route53        = false
  enable_load_balancer  = false
  tags                  = merge(local.tags, var.tags)

  # app specific variables
  name                              = "app"
  image                             = local.image_app
  cpu                               = "512"
  memory                            = "1024"
  environment                       = local.config_app
  secrets                           = local.secrets
  kms_secrets                       = local.kms_secrets_app
  desired_count                     = var.desired_count
  health_check_grace_period_seconds = 300
  enable_execute_command            = true
}
