terraform {
  backend "s3" {}
}

provider "aws" {
  region = "us-east-2"
}

locals {
  git = "gemini"
}

data "aws_route53_zone" "this" {
  name = "oss.champtest.net."
}

data "aws_vpcs" "this" {
  tags = {
    purpose = "vega"
  }
}

data "aws_subnets" "private" {
  tags = {
    purpose = "vega"
    Type    = "Private"
  }

  filter {
    name   = "vpc-id"
    values = [data.aws_vpcs.this.ids[0]]
  }
}

data "aws_subnets" "public" {
  tags = {
    purpose = "vega"
    Type    = "Public"
  }

  filter {
    name   = "vpc-id"
    values = [data.aws_vpcs.this.ids[0]]
  }
}

module "acm" {
  source            = "github.com/champ-oss/terraform-aws-acm.git?ref=v1.0.115-bfc08dd"
  git               = local.git
  domain_name       = "${local.git}.${data.aws_route53_zone.this.name}"
  create_wildcard   = false
  zone_id           = data.aws_route53_zone.this.zone_id
  enable_validation = true
}

module "kms" {
  source                  = "github.com/champ-oss/terraform-aws-kms.git?ref=v1.0.31-3fc28eb"
  git                     = local.git
  name                    = "alias/${local.git}-test"
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
  certificate_arn        = module.acm.arn
  github_app_id          = aws_kms_ciphertext.github_app_id.ciphertext_blob
  github_installation_id = aws_kms_ciphertext.github_installation_id.ciphertext_blob
  github_pem             = aws_kms_ciphertext.github_pem.ciphertext_blob
  private_subnet_ids     = data.aws_subnets.private.ids
  public_subnet_ids      = data.aws_subnets.public.ids
  vpc_id                 = data.aws_vpcs.this.ids[0]
  domain                 = data.aws_route53_zone.this.name
  zone_id                = data.aws_route53_zone.this.zone_id
  protect                = false
  grafana_force_oauth    = false
  use_terraform_api_key  = false
  minutes_between_checks = 0.25
  drop_tables            = true
  repos = [
    "champ-oss/terraform-env-template"
  ]
}
