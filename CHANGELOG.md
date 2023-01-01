## 0.1.0 (27 december, 2022)

NOTES:

This is the first release of the Elestio Terraform Provider.

FEATURES:

- **New Resource:** `elestio_project`
- **New Resource:** `elestio_service`
- **New Datasource:** `elestio_project`

## 0.2.0 (29 december, 2022)

NOTES:

This release fix many bugs and introduce all the available service templates provider by Elestio API.
The documentation was also improved.

FIX:

- Forbid the update of `elestio_postgres`.`version`. It **requires** a replace of the full resource.
- Fix error happening if creating a service with `lightsail` provider.

DEPRECATIONS:

- **Deprecated Resource:** `elestio_postgres`. Use `elestio_postgresql` instead.

FEATURES:

- **New Resource:** for every service templates
- Add default `version` value for template services if there is one recommended by the API.
- Improve documentation

## 0.3.0 (30 december, 2022)

FEATURES:

- **New Resource:** `elestio_couchbase`
- **New Resource:** `elestio_searxng`
- **Improve Documentation:** Add docker hub image link of services
- **New guide:** Get started

## 0.4.0 (Unreleased)

The available templates are now saved in a JSON file in the repo.
This avoids publishing new resources by mistake during a build.

FEATURES:

- **New Resource:** `elestio_manticoresearch`
- **Improved Documentation:** Add examples for each templates
