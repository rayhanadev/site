# -----------------------------------------------------------------------------
# Tailscale Auth Key
# -----------------------------------------------------------------------------

resource "tailscale_tailnet_key" "server" {
  reusable      = true
  preauthorized = true
  tags          = ["tag:server"]
  description   = "Terraform-managed key for ${var.instance_name}"

  depends_on = [tailscale_acl.policy]
}

# -----------------------------------------------------------------------------
# Tailscale ACL Policy
# -----------------------------------------------------------------------------
#
# WARNING: This resource overwrites the entire ACL policy for the tailnet.
# Export your current ACL before the first `terraform apply`:
#   tailscale policy get > acl-backup.json
#

resource "tailscale_acl" "policy" {
  overwrite_existing_content = true

  acl = jsonencode({
    tagOwners = {
      "tag:server" = ["autogroup:admin"]
      "tag:ci"     = ["autogroup:admin"]
    }

    acls = [
      {
        action = "accept"
        src    = ["autogroup:admin"]
        dst    = ["*:*"]
      },
      {
        action = "accept"
        src    = ["tag:ci"]
        dst    = ["tag:server:22"]
      },
    ]

    ssh = [
      {
        action = "accept"
        src    = ["autogroup:admin"]
        dst    = ["tag:server"]
        users  = ["autogroup:nonroot", "root"]
      },
      {
        action = "accept"
        src    = ["tag:ci"]
        dst    = ["tag:server"]
        users  = ["autogroup:nonroot"]
      },
    ]
  })
}
