# Create and manage a load balancer via terraform
resource "elestio_load_balancer" "loadbalancer" {
  project_id    = "project-id"
  provider_name = "scaleway"
  datacenter    = "fr-par-1"
  server_type   = "SMALL-2C-2G"
  config = {
    target_services = [
      // Services ids, IPs or CNAMEs
      elestio_service.service_1.id,
      elestio_service.service_2.ipv4,
      elestio_service.service_3.cname,
      "212.47.xxx.xx",
      "myawesomeapp.com",
    ]
  }
}

resource "elestio_service" "service_1" {
  // ...
}

resource "elestio_service" "service_2" {
  // ...
}

resource "elestio_service" "service_3" {
  // ...
}
