locals {
  config_grafana = {
    GF_DATABASE_HOST                     = aws_rds_cluster.this.endpoint
    GF_DATABASE_TYPE                     = "mysql"
    GF_DATABASE_USER                     = var.database_username
    GF_SECURITY_ADMIN_USER               = var.grafana_username
    GF_SERVER_ROOT_URL                   = "https://${local.grafana_dns}/"
    GF_AUTH_BASIC_ENABLED                = var.grafana_force_oauth ? "false" : "true"
    GF_AUTH_DISABLE_LOGIN_FORM           = var.grafana_force_oauth ? "true" : "false"
    GF_AUTH_OAUTH_AUTO_LOGIN             = var.grafana_force_oauth ? "true" : "false"
    GF_AUTH_GITHUB_ENABLED               = "true"
    GF_AUTH_GITHUB_ALLOW_SIGN_UP         = "true"
    GF_AUTH_GITHUB_SCOPES                = "user:email,read:org"
    GF_AUTH_GITHUB_AUTH_URL              = "https://github.com/login/oauth/authorize"
    GF_AUTH_GITHUB_TOKEN_URL             = "https://github.com/login/oauth/access_token"
    GF_AUTH_GITHUB_API_URL               = "https://api.github.com/user"
    GF_AUTH_GITHUB_TEAM_IDS              = ""
    GF_AUTH_GITHUB_ALLOWED_ORGANIZATIONS = "champtitles"
    GF_AUTH_GITHUB_CLIENT_ID             = var.grafana_github_client_id
  }

  gf_auth_github_client_secret = var.grafana_github_client_secret != "" ? { GF_AUTH_GITHUB_CLIENT_SECRET = var.grafana_github_client_secret } : {}

  grafana_dns = var.grafana_hostname != null ? "${var.grafana_hostname}.${var.domain}" : "grafana-${random_string.identifier.result}.${var.domain}"
}

resource "random_password" "grafana" {
  length  = 32
  special = false

  lifecycle {
    create_before_destroy = true
  }
}

module "grafana" {
  source                = "github.com/champ-oss/terraform-aws-app.git?ref=v1.0.175-10268b5"
  git                   = "${var.git}-${random_string.identifier.result}"
  vpc_id                = var.vpc_id
  subnets               = var.private_subnet_ids
  zone_id               = var.zone_id
  cluster               = module.core.ecs_cluster_name
  security_groups       = [module.core.ecs_app_security_group]
  execution_role_arn    = module.core.execution_ecs_role_arn
  wait_for_steady_state = var.wait_for_steady_state
  enable_route53        = true
  enable_load_balancer  = true
  tags                  = merge(local.tags, var.tags)
  listener_arn          = module.core.lb_public_listener_arn
  lb_dns_name           = module.core.lb_public_dns_name
  lb_zone_id            = module.core.lb_public_zone_id

  # app specific variables
  name                              = "grafana"
  dns_name                          = local.grafana_dns
  image                             = var.image_grafana
  healthcheck                       = "/api/health"
  cpu                               = "512"
  memory                            = "1024"
  environment                       = local.config_grafana
  secrets                           = local.secrets
  kms_secrets                       = merge(local.gf_auth_github_client_secret)
  desired_count                     = var.desired_count_grafana
  health_check_grace_period_seconds = 300
  enable_execute_command            = true
  port                              = 3000
}
