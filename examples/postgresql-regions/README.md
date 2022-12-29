# Deploying PostgreSQL databases in multiple regions

## Requirements

- [Signup for Elestio](https://dash.elest.io/signup)
- [Get your API Token](https://dash.elest.io/account/security)
- [Install terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli#install-terraform)
- (Optional) [Install psql](https://www.timescale.com/blog/how-to-install-psql-on-mac-ubuntu-debian-windows/)
- [Download exemple folder](https://github.com/elestio/terraform-provider-elestio/exemples/postgresql-regions).

```
|- .gitignore
|- README.md
|- main.tf
|- secret.tfvars.tmp
|- variables.tf
```

## Variables

Rename `./secrets.tfvars.tmp` to `./secrets.tfvars` and fill in the appropriate values.

## Initialize Terraform

Ensure that you have Terraform installed and initialize the project.

```sh
$ terraform init

Initializing the backend...

Initializing provider plugins...

Terraform has been successfully initialized!
```

## Plan

You can preview the different resources that terraform will create following the `main.tf`.

```sh
$ terraform plan -var-file="secret.tfvars"

Plan: 3 to add, 0 to change, 0 to destroy.
```

## Apply

Deploy your infrastructure.

```sh
$ terraform apply -var-file="secret.tfvars"

Plan: 3 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

elestio_project.pg_project: Creating...
elestio_project.pg_project: Creation complete after 0s [id=2318]
elestio_postgresql.pg_asia: Creating...
elestio_postgresql.pg_europe: Creating...
elestio_postgresql.pg_asia: Still creating... [10s elapsed]
elestio_postgresql.pg_europe: Still creating... [10s elapsed]
...
elestio_postgresql.pg_asia: Creation complete after 5m0s [id=pg-asia-u525.vm.elestio.app]
elestio_postgresql.pg_europe: Creation complete after 5m0s [id=pg-europe-u525.vm.elestio.app]


Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
```

## (Optional) Access the databases

- Run these commands in your terminal (you need to install `psql`):

```bash
eval "$(terraform output -raw pg_europe_psql)"
eval "$(terraform output -raw pg_asia_psql)"
```

## Cleanup

```sh
$ terraform destroy -var-file="secret.tfvars"


Plan: 0 to add, 0 to change, 3 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

elestio_postgresql.pg_asia: Destroying... [id=pg-asia-u525.vm.elestio.app]
elestio_postgresql.pg_europe: Destroying... [id=pg-europe-u525.vm.elestio.app]
elestio_postgresql.pg_asia: Destruction complete after 0s
elestio_postgresql.pg_europe: Destruction complete after 0s
elestio_project.pg_project: Destroying... [id=2318]
elestio_project.pg_project: Destruction complete after 0s


Destroy complete! Resources: 3 destroyed.
```
