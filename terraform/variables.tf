variable "repos" {
  description = "List of Github repository names to gather data"
  type        = list(string)
}

variable "zone_id" {
  description = "Route53 Zone ID"
  type        = string
}

variable "domain" {
  description = "Route53 Domain"
  type        = string
}

variable "vpc_id" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group#vpc_id"
  type        = string
}

variable "public_subnet_ids" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb#subnets"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster#subnet_ids"
  type        = list(string)
}

variable "certificate_arn" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener#certificate_arn"
  type        = string
}

variable "github_app_id" {
  description = "KMS encrypted GitHub App ID"
  type        = string
}

variable "github_installation_id" {
  description = "KMS encrypted GitHub Installation ID"
  type        = string
}

variable "github_pem" {
  description = "KMS encrypted GitHub Private Key PEM which is base64 encoded"
  type        = string
}

variable "grafana_github_client_id" {
  description = "https://github.com/organizations/champtitles/settings/applications"
  type        = string
  default     = ""
}

variable "grafana_github_client_secret" {
  description = "https://github.com/organizations/champtitles/settings/applications (KMS encrypted)"
  type        = string
  default     = ""
}

variable "additional_certificate_arns" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener_certificate#certificate_arn"
  type        = list(string)
  default     = []
}

variable "wait_for_steady_state" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service#wait_for_steady_state"
  type        = bool
  default     = true
}

variable "desired_count" {
  description = "https://www.terraform.io/docs/providers/aws/r/ecs_service.html#desired_count"
  type        = number
  default     = 1
}

variable "desired_count_grafana" {
  description = "https://www.terraform.io/docs/providers/aws/r/ecs_service.html#desired_count"
  type        = number
  default     = 1
}

variable "database_max_capacity" {
  description = "https://www.terraform.io/docs/providers/aws/r/rds_cluster.html#max_capacity"
  type        = number
  default     = 8
}

variable "protect" {
  description = "Enables deletion protection on eligible resources"
  type        = bool
  default     = true
}

variable "snapshot_identifier" {
  description = "https://www.terraform.io/docs/providers/aws/r/rds_cluster.html#snapshot_identifier"
  type        = string
  default     = ""
}

variable "db_cluster_parameter_group_name" {
  description = "https://www.terraform.io/docs/providers/aws/r/rds_cluster.html#db_cluster_parameter_group_name"
  type        = string
  default     = ""
}

variable "database_username" {
  description = "Default database username"
  type        = string
  default     = "root"
}

variable "image_grafana" {
  description = "Docker image for Grafana"
  type        = string
  default     = "grafana/grafana:8.3.4"
}

variable "grafana_username" {
  description = "Default Grafana username"
  type        = string
  default     = "admin"
}

variable "tags" {
  description = "Map of tags to assign to resources"
  type        = map(string)
  default     = {}
}

variable "git" {
  description = "Name of the Git repo"
  type        = string
  default     = "gemini"
}

variable "grafana_force_oauth" {
  description = "Disable basic auth and force users to sign in with OAuth. Enable this once OAuth is tested and working."
  type        = bool
  default     = false
}

variable "grafana_hostname" {
  description = "Optional hostname for the Grafana server. If omitted a random identifier will be used."
  type        = string
  default     = "gemini"
}

variable "use_terraform_api_key" {
  description = "Use API key to authenticate with the Grafana Terraform Provider instead of basic auth"
  type        = bool
  default     = false
}

variable "terraform_api_key" {
  description = "KMS encrypted API key (generated manually from the Grafana UI"
  type        = string
  default     = ""
}

variable "debug" {
  description = "Enable Gemini debug logging"
  type        = bool
  default     = true
}

variable "minutes_between_checks" {
  description = "How often Gemini should poll for Github data"
  type        = number
  default     = 5
}

variable "cluster_instance_count" {
  description = "Database cluster instances"
  type        = number
  default     = 1
}

variable "metric_alarms_enabled" {
  description = "enable or disable cloudwatch metric alarms for aurora"
  type        = bool
  default     = true
}