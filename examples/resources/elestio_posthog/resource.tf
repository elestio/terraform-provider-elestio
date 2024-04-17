resource "elestio_posthog" "example" {
  project_id    = "2500"
  version       = "92e17ce307a577c4233d4ab252eebc6c2207a5ee"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
