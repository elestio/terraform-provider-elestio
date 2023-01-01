# Create and manage iTop service.
resource "elestio_itop" "my_itop" {
  project_id    = "2500"
  server_name   = "awesome-itop"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
