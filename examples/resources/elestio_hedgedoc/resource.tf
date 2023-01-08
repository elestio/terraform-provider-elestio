# Create and manage HedgeDoc service.
resource "elestio_hedgedoc" "my_hedgedoc" {
  project_id    = "2500"
  server_name   = "awesome-hedgedoc"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
