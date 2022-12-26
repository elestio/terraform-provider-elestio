# Create and manage a project.
resource "elestio_project" "project1" {
  name             = "My custom name"
  description      = "My fancy description"
  technical_emails = "admin@exemple.com"
}
