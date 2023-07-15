# Create and manage Appwrite service.
resource "elestio_appwrite" "demo_appwrite" {
  project_id    = "2500"
  server_name   = "demo-appwrite"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
