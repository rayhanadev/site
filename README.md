# rayhanadev.com

Go HTTP/3 server that streams content character-by-character. Serves HTML to browsers and ANSI-styled text to terminals (`curl https://rayhanadev.com`).

## Development

Requires Go 1.25+ and TLS certificates at `configs/certs/{cert,key}.pem`.

```sh
task dev          # run dev server on :3000
task build        # compile to bin/site
task lint         # golangci-lint
task format       # gofmt
```

## Deployment

Pushes to `main` trigger three GitHub Actions jobs:

1. **Build image** — builds the Docker image and pushes it to `ghcr.io/rayhanadev/site`
2. **Deploy server** (after build) — SSHs into the GCP instance via Tailscale, pulls the new image, and restarts via Docker Compose
3. **Deploy static assets** — deploys `public/` to Cloudflare Workers via Wrangler

### GitHub Actions Secrets

| Secret | Description |
|--------|-------------|
| `TS_OAUTH_CLIENT_ID` | Tailscale OAuth client ID |
| `TS_OAUTH_SECRET` | Tailscale OAuth client secret |
| `CLOUDFLARE_API_TOKEN` | Cloudflare API token for Wrangler |
| `CLOUDFLARE_ACCOUNT_ID` | Cloudflare account ID |

## Infrastructure (Terraform)

All infrastructure is codified in `terraform/`. Managed resources:

- **GCP** — VPC, subnet, static IP, firewall (Cloudflare-only origin access), e2-micro instance with cloud-init
- **Cloudflare** — DNS records, SSL Full (Strict), Origin CA certificate, www-to-apex redirect
- **Tailscale** — auth key, ACL policy (SSH access for admin and CI)

### Prerequisites

- [Terraform](https://developer.hashicorp.com/terraform/install) or [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.7
- GCP project with billing enabled
- Cloudflare API token (Zone:Edit, DNS:Edit, Origin CA permissions)
- Tailscale OAuth client credentials

### Bootstrap (one-time)

Create the GCS bucket for remote state:

```sh
cd terraform/bootstrap
cp terraform.tfvars.example terraform.tfvars  # fill in gcp_project_id
tofu init
tofu apply
```

### Provision Infrastructure

```sh
cd terraform
cp terraform.tfvars.example terraform.tfvars  # fill in all values

# Update the backend bucket name in terraform.tf to match your project ID,
# or pass it at init time:
tofu init -backend-config="bucket=YOUR_PROJECT_ID-terraform-state"

tofu plan    # review changes
tofu apply   # provision everything
```
