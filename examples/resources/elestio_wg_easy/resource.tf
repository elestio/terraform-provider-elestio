# Create and manage WG-Easy service.
resource "elestio_wg_easy" "my_wg_easy" {
  project_id    = "2500"
  server_name   = "awesome-wg_easy"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
