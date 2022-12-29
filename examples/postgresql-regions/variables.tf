variable "elestio_email" {
  description = "Elestio Email"
  type        = string
}

variable "elestio_api_token" {
  description = "Elestio API Token"
  type        = string
  sensitive   = true
}

