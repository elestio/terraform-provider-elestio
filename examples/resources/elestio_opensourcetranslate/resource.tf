# Create and manage OpenSourceTranslate service.
resource "elestio_opensourcetranslate" "demo_opensourcetranslate" {
  project_id    = "2500"
  server_name   = "demo-opensourcetranslate"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
