# -----------------------------------------------------------------------------
# Origin Server TLS Certificate (Cloudflare Origin CA)
# -----------------------------------------------------------------------------

# ECDSA P-256 private key for the origin server
resource "tls_private_key" "origin" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P256"
}

# Certificate Signing Request for the domain + wildcard
resource "tls_cert_request" "origin" {
  private_key_pem = tls_private_key.origin.private_key_pem

  subject {
    common_name = var.domain
  }

  dns_names = [
    var.domain,
    "*.${var.domain}",
  ]
}

# Cloudflare Origin CA certificate — 15-year validity
resource "cloudflare_origin_ca_certificate" "origin" {
  csr                = tls_cert_request.origin.cert_request_pem
  hostnames          = ["*.${var.domain}", var.domain]
  request_type       = "origin-ecc"
  requested_validity = 5475
}
