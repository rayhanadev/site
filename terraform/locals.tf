locals {
  labels = {
    app         = var.instance_name
    managed_by  = "terraform"
    environment = "production"
  }
}
