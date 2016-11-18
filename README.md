# docker-machine-driver-terraform

A driver for Docker Machine that uses Terraform to provision infrastructure.

The driver can work with a Terraform configuration in any of the following forms:

* A single local .tf file
* A single local .zip file containing 1 or more .tf files
* A single remote .tf file (using HTTP)
* A single remote .zip file containing 1 or more .tf files (using HTTP)
* A local directory containing 1 or more .tf files
* A directory from a GitHub repository

It will supply the following values to the configuration as variables (in addition to any supplied via `--terraform-variables-json` or `--terraform-variables-json-file`):

* `dm_machine_name` - The name of the Docker machine being created
* `dm_onetime_password` - A one-time password that can be used for scenarios such as bootstrapping key-based SSH authentication

It expects the following [outputs](https://www.terraform.io/docs/configuration/outputs.html) from Terraform:

* `dm_machine_ip` (Required) - The IP address of the target machine
* `dm_machine_ssh_username` (Optional) - The SSH user name for authentication to the target machine

## Installing the driver

Download the [latest release](https://github.com/tintoy/docker-machine-driver-terraform/releases) and place the provider executable in the same directory as `docker-machine` executable (or somewhere on your `PATH`).

## Building the driver

If you'd rather run from source, simply run `make install` and you're good to go.
