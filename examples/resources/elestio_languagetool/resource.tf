# Create and manage LanguageTool service.
resource "elestio_languagetool" "demo_languagetool" {
  project_id    = "2500"
  server_name   = "demo-languagetool"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
