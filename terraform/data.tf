# Cloudflare edge IP ranges — used to restrict origin firewall to Cloudflare only
data "cloudflare_ip_ranges" "cloudflare" {}
