# -----------------------------------------------------------------------------
# VPC and Subnet
# -----------------------------------------------------------------------------

resource "google_compute_network" "vpc" {
  name                    = "${var.instance_name}-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name                     = "${var.instance_name}-subnet"
  ip_cidr_range            = "10.0.1.0/24"
  network                  = google_compute_network.vpc.id
  private_ip_google_access = true
}

# -----------------------------------------------------------------------------
# Static External IP
# -----------------------------------------------------------------------------

resource "google_compute_address" "static_ip" {
  name = "${var.instance_name}-ip"
}

# -----------------------------------------------------------------------------
# Firewall Rules
# -----------------------------------------------------------------------------

# Allow HTTPS (port 8443) from Cloudflare IPv4 edge IPs
resource "google_compute_firewall" "cloudflare_https_ipv4" {
  name    = "${var.instance_name}-allow-cf-https-ipv4"
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["8443"]
  }

  allow {
    protocol = "udp"
    ports    = ["8443"]
  }

  source_ranges = data.cloudflare_ip_ranges.cloudflare.ipv4_cidrs
  target_tags   = ["web-server"]
}

# Allow HTTPS (port 8443) from Cloudflare IPv6 edge IPs
resource "google_compute_firewall" "cloudflare_https_ipv6" {
  name    = "${var.instance_name}-allow-cf-https-ipv6"
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["8443"]
  }

  allow {
    protocol = "udp"
    ports    = ["8443"]
  }

  source_ranges = data.cloudflare_ip_ranges.cloudflare.ipv6_cidrs
  target_tags   = ["web-server"]
}

# Allow SSH from IAP (emergency fallback — Tailscale is primary access)
resource "google_compute_firewall" "iap_ssh" {
  name    = "${var.instance_name}-allow-iap-ssh"
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["35.235.240.0/20"]
  target_tags   = ["ssh-iap"]
}
