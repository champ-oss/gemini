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