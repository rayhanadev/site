output "instance_ip" {
  description = "Static external IP of the GCE instance"
  value       = google_compute_address.static_ip.address
}

output "instance_name" {
  description = "Name of the GCE instance"
  value       = google_compute_instance.server.name
}

output "tailscale_auth_key" {
  description = "Tailscale auth key for manual use if needed"
  value       = tailscale_tailnet_key.server.key
  sensitive   = true
}

output "origin_cert_expires_on" {
  description = "Expiry date of the Cloudflare Origin CA certificate"
  value       = cloudflare_origin_ca_certificate.origin.expires_on
}

output "dns_record_apex_id" {
  description = "Cloudflare DNS record ID for the apex domain"
  value       = cloudflare_dns_record.apex.id
}
