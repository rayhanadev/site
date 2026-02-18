# -----------------------------------------------------------------------------
# Required variables
# -----------------------------------------------------------------------------

variable "gcp_project_id" {
  description = "GCP project ID"
  type        = string
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token with Zone:Edit, DNS:Edit, Origin CA permissions"
  type        = string
  sensitive   = true
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID for rayhanadev.com"
  type        = string
}

variable "tailscale_oauth_client_id" {
  description = "Tailscale OAuth client ID (required for tag-based auth keys)"
  type        = string
  sensitive   = true
}

variable "tailscale_oauth_client_secret" {
  description = "Tailscale OAuth client secret"
  type        = string
  sensitive   = true
}

variable "tailscale_tailnet" {
  description = "Tailscale tailnet name"
  type        = string
}

# -----------------------------------------------------------------------------
# Optional variables with defaults
# -----------------------------------------------------------------------------

variable "gcp_region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "gcp_zone" {
  description = "GCP zone"
  type        = string
  default     = "us-central1-a"
}

variable "domain" {
  description = "Primary domain name"
  type        = string
  default     = "rayhanadev.com"
}

variable "instance_name" {
  description = "Name for the GCE instance"
  type        = string
  default     = "rayhanadevcom"
}

variable "machine_type" {
  description = "GCE machine type (must be e2-micro for free tier)"
  type        = string
  default     = "e2-micro"

  validation {
    condition     = var.machine_type == "e2-micro"
    error_message = "Machine type must be e2-micro to stay within the GCP free tier."
  }
}

variable "site_repo_url" {
  description = "Git repository URL for the site source code"
  type        = string
  default     = "https://github.com/rayhanadev/site.git"
}
