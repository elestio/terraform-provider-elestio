# Create and manage Jupyter Notebook service.
resource "elestio_jupyter" "my_jupyter" {
  project_id    = "2500"
  server_name   = "awesome-jupyter"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
