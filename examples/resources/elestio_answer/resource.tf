# Create and manage Answer service.
resource "elestio_answer" "demo_answer" {
  project_id    = "2500"
  server_name   = "demo-answer"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
