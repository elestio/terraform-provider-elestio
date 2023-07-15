# Create and manage Supabase service.
resource "elestio_supabase" "demo_supabase" {
  project_id    = "2500"
  server_name   = "demo-supabase"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
