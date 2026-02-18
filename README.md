# rayhanadev.com

Personal site — a Go HTTP/3 server that streams content character-by-character with a typewriter effect. Serves HTML to browsers and ANSI-styled text to terminals (`curl https://rayhanadev.com`).

## Architecture

```
                    ┌──────────────────────────┐
                    │     Cloudflare Edge       │
                    │                           │
  User ──HTTPS──►  │  CDN (static assets)      │
                    │  DNS (rayhanadev.com)      │
                    │  SSL termination           │
                    │  www→apex redirect         │
                    │                           │
                    └──────────┬───────────────┘
                               │ Origin (port 8443, Full Strict)
                               ▼
                    ┌──────────────────────────┐
                    │  GCP e2-micro (free tier) │
                    │                           │
                    │  Docker → Go HTTP/3 server│
                    │  Cloudflare Origin CA cert │
                    │  Tailscale SSH enabled     │
                    └──────────────────────────┘
                               ▲
                               │ Tailscale tunnel
                    ┌──────────┴───────────────┐
                    │  GitHub Actions (deploy)  │
                    │  SSH via Tailscale → pull  │
                    │  + docker compose up      │
                    └──────────────────────────┘
```

## Local Development

Requires Go 1.25+ and TLS certificates at `configs/certs/{cert,key}.pem`.

```sh
task dev          # run dev server on :3000
task build        # compile to bin/site
task lint         # golangci-lint
task format       # gofmt
```

## Deployment

Pushes to `main` trigger two parallel GitHub Actions jobs:

1. **Go server** — connects to Tailscale, SSHs into the GCP instance, pulls the latest code, and rebuilds via Docker Compose
2. **Static assets** — deploys `public/` to Cloudflare Workers via Wrangler

### GitHub Actions Secrets

| Secret | Description |
|--------|-------------|
| `TS_OAUTH_CLIENT_ID` | Tailscale OAuth client ID |
| `TS_OAUTH_SECRET` | Tailscale OAuth client secret |
| `DEPLOY_HOST` | Tailscale hostname of the server |
| `DEPLOY_USER` | SSH user on the server |
| `DEPLOY_SSH_KEY` | SSH private key for deployment |
| `CLOUDFLARE_API_TOKEN` | Cloudflare API token for Wrangler |

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

### Terraform File Structure

```
terraform/
├── terraform.tf            # version constraints, backend
├── providers.tf            # google, cloudflare, tailscale, tls
├── variables.tf            # input variables
├── locals.tf               # derived values
├── data.tf                 # cloudflare_ip_ranges
├── networking.tf           # VPC, subnet, static IP, firewall
├── compute.tf              # service account, GCE instance
├── tls.tf                  # ECDSA key, CSR, Origin CA cert
├── dns.tf                  # DNS records, SSL settings, redirect
├── tailscale.tf            # auth key, ACL policy
├── outputs.tf              # instance IP, cert expiry, etc.
├── templates/
│   └── cloud-init.yaml.tftpl
├── bootstrap/
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
└── terraform.tfvars.example
```

### Important Notes

- **Tailscale ACL**: `tailscale_acl` overwrites the entire tailnet policy. Back up your current ACL before the first apply: `tailscale policy get > acl-backup.json`
- **Cloud-init**: Runs once on first boot. Changes to the template force instance recreation (immutable infrastructure).
- **Origin CA cert**: 15-year Cloudflare Origin CA certificate. Private key lives in Terraform state — keep the GCS backend encrypted.
- **Firewall**: The origin server only accepts traffic from Cloudflare edge IPs on port 8443. IAP SSH on port 22 is available as an emergency fallback.
