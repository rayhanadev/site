<img alt="image" src="https://github.com/user-attachments/assets/d5d6a6d6-a43f-48ce-b4ef-6cae4feab856" />

# rayhanadev[dot]com

Go HTTP/3 server with a custom typewriter markup parser that streams content character-by-character to browsers and renders ANSI-styled text for terminals (`curl https://rayhanadev.com`).

## Development

Requires Go 1.25+ and TLS certificates (defaults to `configs/certs/{cert,key}.pem`, override with `TLS_CERT_PATH` / `TLS_KEY_PATH` env vars). Listen address defaults to `:3000` (override with `LISTEN_ADDR`).

```sh
task dev                  # run dev server on :3000
task build                # compile to bin/site
task lint                 # golangci-lint
task format               # gofmt
task generate:wrangler    # regenerate wrangler.jsonc from static assets
task deploy:static        # generate wrangler config + deploy to Cloudflare Workers
task deploy:server:build  # build and push Docker image to GHCR
task deploy:server        # SSH into production and pull latest image
```

## Deployment

Two path-filtered GitHub Actions workflows run on pushes to `main`:

**[Deploy Server](.github/workflows/deploy.yml)** — triggered by changes to Go source, `Dockerfile`, or `docker-compose.yml` (excludes `internal/assets/static/`):

1. **Build image** — builds the Docker image and pushes it to `ghcr.io/rayhanadev/site`
2. **Deploy server** (after build) — SSHs into the GCP instance via Tailscale, pulls the new image, and restarts via Docker Compose

**[Deploy Static Assets](.github/workflows/deploy-static.yml)** — triggered by changes to `internal/assets/static/` or the wrangler config generator:

1. **Deploy static assets** — generates `wrangler.jsonc` from the static assets directory and deploys to Cloudflare Workers via Wrangler

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
