# -----------------------------------------------------------------------------
# Service Account
# -----------------------------------------------------------------------------

resource "google_service_account" "instance" {
  account_id   = "${var.instance_name}-sa"
  display_name = "Personal site instance service account"
}

# -----------------------------------------------------------------------------
# Compute Instance
# -----------------------------------------------------------------------------

resource "google_compute_instance" "server" {
  name         = var.instance_name
  machine_type = var.machine_type
  tags         = ["web-server", "ssh-iap"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2404-lts-amd64"
      size  = 30
      type  = "pd-standard"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnet.id

    access_config {
      nat_ip = google_compute_address.static_ip.address
    }
  }

  metadata = {
    user-data = templatefile("${path.module}/templates/cloud-init.yaml.tftpl", {
      tailscale_auth_key = tailscale_tailnet_key.server.key
      origin_cert_pem    = cloudflare_origin_ca_certificate.origin.certificate
      origin_key_pem     = tls_private_key.origin.private_key_pem
      site_repo_url      = var.site_repo_url
    })
  }

  service_account {
    email  = google_service_account.instance.email
    scopes = ["cloud-platform"]
  }

  shielded_instance_config {
    enable_secure_boot          = true
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }

  scheduling {
    preemptible       = false
    automatic_restart = true
  }

  labels = local.labels

  allow_stopping_for_update = true
}
