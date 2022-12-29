## 0.1.0 (27 december, 2022)

NOTES:

This is the first release of the Elestio Terraform Provider.

FEATURES:

- **New Resource:** `elestio_project`
- **New Resource:** `elestio_service`
- **New Datasource:** `elestio_project`

## 0.1.1 (29 december, 2022)

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
