# Create and manage Answer service.
resource "elestio_answer" "my_answer" {
  project_id    = "2500"
  server_name   = "awesome-answer"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
