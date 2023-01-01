# Create and manage Huginn service.
resource "elestio_huginn" "my_huginn" {
  project_id    = "2500"
  server_name   = "awesome-huginn"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
