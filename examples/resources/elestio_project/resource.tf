# Create and manage a project.
resource "elestio_project" "myawesomeproject" {
  name             = "Awesome project"
  technical_emails = "YOUR-EMAIL"
}
