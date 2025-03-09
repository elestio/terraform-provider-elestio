resource "elestio_supabase" "example" {
  project_id    = "2500"
  version       = "20241202-71e5240"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
