# -----------------------------------------------------------------------------
# DNS Records
# -----------------------------------------------------------------------------

# Apex A record — proxied through Cloudflare
resource "cloudflare_dns_record" "apex" {
  zone_id = var.cloudflare_zone_id
  name    = "@"
  type    = "A"
  content = google_compute_address.static_ip.address
  proxied = true
  ttl     = 1 # Auto when proxied
}

# www A record — proxied, stub target for redirect rule
resource "cloudflare_dns_record" "www" {
  zone_id = var.cloudflare_zone_id
  name    = "www"
  type    = "A"
  content = google_compute_address.static_ip.address
  proxied = true
  ttl     = 1
}

# -----------------------------------------------------------------------------
# SSL/TLS Settings
# -----------------------------------------------------------------------------

resource "cloudflare_zone_setting" "ssl" {
  zone_id    = var.cloudflare_zone_id
  setting_id = "ssl"
  value      = "strict"
}

resource "cloudflare_zone_setting" "always_use_https" {
  zone_id    = var.cloudflare_zone_id
  setting_id = "always_use_https"
  value      = "on"
}

resource "cloudflare_zone_setting" "min_tls_version" {
  zone_id    = var.cloudflare_zone_id
  setting_id = "min_tls_version"
  value      = "1.2"
}

# -----------------------------------------------------------------------------
# www → apex redirect
# -----------------------------------------------------------------------------

resource "cloudflare_ruleset" "www_redirect" {
  zone_id     = var.cloudflare_zone_id
  name        = "www-redirect"
  description = "Redirect www.${var.domain} to ${var.domain}"
  kind        = "zone"
  phase       = "http_request_dynamic_redirect"

  rules = [{
    ref         = "www_to_apex"
    description = "301 redirect www to apex"
    expression  = "(http.host eq \"www.${var.domain}\")"
    action      = "redirect"
    action_parameters = {
      from_value = {
        status_code = 301
        target_url = {
          expression = "concat(\"https://${var.domain}\", http.request.uri.path)"
        }
        preserve_query_string = true
      }
    }
  }]
}
