# Create and manage Rstudio service.
resource "elestio_rstudio" "my_rstudio" {
  project_id    = "2500"
  server_name   = "awesome-rstudio"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
