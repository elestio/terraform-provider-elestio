# Create and manage Supabase service.
resource "elestio_supabase" "my_supabase" {
  project_id    = "2500"
  server_name   = "awesome-supabase"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
