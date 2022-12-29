provider "elestio" {
  email     = "YOUR-EMAIL"
  api_token = "YOUR-API-TOKEN"
}

resource "elestio_project" "myproject" {
  # ...
}
