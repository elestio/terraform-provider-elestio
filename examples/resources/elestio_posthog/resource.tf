# Create and manage PostHog service.
resource "elestio_posthog" "demo_posthog" {
  project_id    = "2500"
  server_name   = "demo-posthog"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
