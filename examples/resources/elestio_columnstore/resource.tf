# Create and manage ColumnStore service.
resource "elestio_columnstore" "demo_columnstore" {
  project_id    = "2500"
  server_name   = "demo-columnstore"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
