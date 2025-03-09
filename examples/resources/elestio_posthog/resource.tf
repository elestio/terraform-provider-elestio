resource "elestio_posthog" "example" {
  project_id    = "2500"
  version       = "389a8d4daa1953c7208ce0b201555e0fbe41674b"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
