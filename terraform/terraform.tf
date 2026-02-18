terraform {
  required_version = "~> 1.7"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
    tailscale = {
      source  = "tailscale/tailscale"
      version = "~> 0.17"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }

  # NOTE: Update bucket name to match your GCP project ID.
  # Must match the bucket created by bootstrap/main.tf: {project_id}-terraform-state
  # Alternatively, pass at init time: terraform init -backend-config="bucket=YOUR_PROJECT-terraform-state"
  backend "gcs" {
    bucket = "rayhanadevcom-terraform-state"
    prefix = "rayhanadevcom"
  }
}
