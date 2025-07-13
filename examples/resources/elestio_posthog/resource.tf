resource "elestio_posthog" "example" {
  project_id    = "2500"
  version       = "8ae62abb31f38eed95d24abd03f100f4cf1d52a1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
