resource "elestio_supabase" "example" {
  project_id    = "2500"
  version       = "20240422-5cf8f30"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
