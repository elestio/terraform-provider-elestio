terraform {
  required_providers {
    elestio = {
      source  = "elestio/elestio"
      version = "0.2.0" # check out the latest version in the release section
    }
  }
}

# Configure the provider
provider "elestio" {
  email     = "<elestio_email>"
  api_token = "<elestio_api_token>"
}

# Create a project
resource "elestio_project" "myawesomeproject" {
  name = "Awesome project"
  # ...
}

