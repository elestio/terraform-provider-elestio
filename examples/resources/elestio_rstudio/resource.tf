# Create and manage Rstudio service.
resource "elestio_rstudio" "demo_rstudio" {
  project_id    = "2500"
  server_name   = "demo-rstudio"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
