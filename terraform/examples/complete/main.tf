provider "aws" {
  region = "us-west-1"
}

locals {
  git = "gemini"
}

data "aws_route53_zone" "this" {
  name = "oss.champtest.net."
}

module "vpc" {
  source                   = "github.com/champ-oss/terraform-aws-vpc.git?ref=v1.0.1-afc8890"
  git                      = local.git
  availability_zones_count = 2
  retention_in_days        = 1
}

module "acm" {
  source            = "github.com/champ-oss/terraform-aws-acm.git?ref=v1.0.1-1cb7679"
  git               = local.git
  domain_name       = data.aws_route53_zone.this.name
  zone_id           = data.aws_route53_zone.this.zone_id
  enable_validation = true
}

module "kms" {
  source                  = "github.com/champ-oss/terraform-aws-kms.git?ref=v1.0.12-2dcefbc"
  git                     = local.git
  name                    = "alias/${local.git}"
  deletion_window_in_days = 7
  account_actions         = []
}

resource "aws_kms_ciphertext" "github_app_id" {
  key_id    = module.kms.key_id
  plaintext = var.github_app_id
}

resource "aws_kms_ciphertext" "github_installation_id" {
  key_id    = module.kms.key_id
  plaintext = var.github_installation_id
}

resource "aws_kms_ciphertext" "github_pem" {
  key_id    = module.kms.key_id
  plaintext = var.github_pem
}

module "this" {
  source                 = "../../"
  docker_tag             = var.docker_tag
  certificate_arn        = module.acm.arn
  github_app_id          = aws_kms_ciphertext.github_app_id.ciphertext_blob
  github_installation_id = aws_kms_ciphertext.github_installation_id.ciphertext_blob
  github_pem             = aws_kms_ciphertext.github_pem.ciphertext_blob
  private_subnet_ids     = module.vpc.private_subnets_ids
  public_subnet_ids      = module.vpc.public_subnets_ids
  vpc_id                 = module.vpc.vpc_id
  domain                 = data.aws_route53_zone.this.name
  zone_id                = data.aws_route53_zone.this.zone_id
  repos = [
    "champtitles/tflint-ruleset-champtitles"
  ]
}