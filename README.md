# docker-machine-driver-terraform

Do you use Docker Machine?
Do you sometimes wish you had a little more control over the infrastructure it creates?

Well now you do. If you can express it as a Terraform configuration, you can use it from Docker Machine.

### Arguments

The driver accepts the following arguments:

* `--terraform-config` (Required) - The path (or URL) of the Terraform configuration to use
* `--terraform-additional-variables` (Optional) - An optional file containing the JSON that represents additional variables for the Terraform configuration
* `--terraform-refresh` (Optional) - A flag which, if specified, will cause the driver to refresh the configuration after applying it

### Terraform configuration

The driver can work with a Terraform configuration in any of the following formats:

* A single local .tf file
* A single local .zip file containing 1 or more .tf files
* A local directory containing 1 or more .tf files
* A single remote .tf file (using HTTP)
* A single remote .zip file containing 1 or more .tf files (using HTTP)

It will supply the following values to the configuration as variables (in addition to any supplied via `--terraform-variables-file`):

* `dm_client_ip` - The public IP of the client machine (useful for configuring firewall rules)
* `dm_machine_name` - The name of the Docker machine being created
* `dm_ssh_user` - The SSH user name to use for authentication
* `dm_ssh_port` - The SSH port to use
* `dm_ssh_public_key_file` - The public SSH key file to use for authentication
* `dm_ssh_private_key_file` - The private SSH key file to use for authentication
* `dm_onetime_password` - An optional one-time password that can be used for scenarios such as bootstrapping key-based SSH authentication

It expects the following [outputs](https://www.terraform.io/docs/configuration/outputs.html) from Terraform:

* `dm_machine_ip` (Required) - The IP address of the target machine
* `dm_machine_ssh_username` (Optional) - The SSH user name for authentication to the target machine  
If specified this overrides the variable of the same name that was passed in

#### Examples

Ok, so it's almost 1am. I'll add some sample TF configurations tomorrow. For now, here's the one I've been using:

```hcl
/*
 * A simple configuration for docker-machine-driver-terraform
 * ----------------------------------------------------------
 */

# Docker Machine variables (supplied by the driver)
variable "dm_client_ip"           { }
variable "dm_machine_name"        { }
variable "dm_ssh_user"            { }
variable "dm_ssh_port"            { }
variable "dm_ssh_public_key_file" { }

# Additional variables (supplied via tfvars.json)
variable "region"                 { }
variable "networkdomain"          { }
variable "datacenter"             { }
variable "vlan"                   { }
variable "ssh_bootstrap_password" { }

# CloudControl
provider "ddcloud" {
  region    = "${var.region}"
}

# Look up network and VLAN.
data "ddcloud_networkdomain" "docker_machine" {
  name        = "${var.networkdomain}"
  datacenter  = "${var.datacenter}"
}
data "ddcloud_vlan" "docker_machine" {
  name          = "${var.vlan}"
  networkdomain = "${data.ddcloud_networkdomain.docker_machine.id}"
}

# Server
resource "ddcloud_server" "docker_machine" {
  name                  = "${var.dm_machine_name}"
  description           = "${var.dm_machine_name} (created by Docker Machine)."
  admin_password        = "${var.ssh_bootstrap_password}"
  auto_start            = true

  memory_gb             = 8

  networkdomain         = "${data.ddcloud_networkdomain.docker_machine.id}"
  primary_adapter_vlan  = "${data.ddcloud_vlan.docker_machine.id}"
  dns_primary           = "8.8.8.8"
  dns_secondary         = "8.8.4.4"

  os_image_name         = "Ubuntu 14.04 2 CPU"

  disk {
    scsi_unit_id        = 0
    size_gb             = 10
  }

  tag {
    name  = "role"
    value = "tf-test"
  }
}

# Server exposure
resource "ddcloud_nat" "docker_machine" {
  networkdomain = "${data.ddcloud_networkdomain.docker_machine.id}"
  private_ipv4  = "${ddcloud_server.docker_machine.primary_adapter_ipv4}"
}
resource "ddcloud_firewall_rule" "docker_machine_ssh4_in" {
  name                = "${replace(var.dm_machine_name, "-", ".")}.ssh4.inbound"
  placement           = "first"
  action              = "accept"
  enabled             = true

  ip_version          = "ipv4"
  protocol            = "tcp"

  source_address      = "${var.dm_client_ip}"

  destination_address = "${ddcloud_nat.docker_machine.public_ipv4}"
  destination_port    = "22"

  networkdomain       = "${data.ddcloud_networkdomain.docker_machine.id}"
}
resource "ddcloud_firewall_rule" "docker_machine_docker_in" {
  name                = "${replace(var.dm_machine_name, "-", ".")}.docker.inbound"
  placement           = "first"
  action              = "accept"
  enabled             = true

  ip_version          = "ipv4"
  protocol            = "tcp"

  source_address      = "${var.dm_client_ip}"

  destination_address = "${ddcloud_nat.docker_machine.public_ipv4}"
  destination_port    = "2376"

  networkdomain       = "${data.ddcloud_networkdomain.docker_machine.id}"

  depends_on          = [ "ddcloud_firewall_rule.docker_machine_ssh4_in" ]
}

# Server SSH bootstrap
#
# Install the SSH key expected by Docker Machine
resource "null_resource" "docker_machine_ssh" {
  # Install our SSH public key.
  provisioner "remote-exec" {
    inline = [
      "mkdir -p ~/.ssh",
      "chmod 700 ~/.ssh",
      "echo '${file(var.dm_ssh_public_key_file)}' > ~/.ssh/authorized_keys",
      "chmod 600 ~/.ssh/authorized_keys",
      "passwd -d root"
    ]

    connection {
      type      = "ssh"

      user      = "root"
      password  = "${var.ssh_bootstrap_password}"

      host      = "${ddcloud_nat.docker_machine.public_ipv4}"
    }
  }

  depends_on    = [ "ddcloud_firewall_rule.docker_machine_ssh4_in" ]
}

# Outputs for Docker Machine
output "dm_machine_ip" {
  value = "${ddcloud_nat.docker_machine.public_ipv4}"
}
output "dm_ssh_user" {
  value = "${var.dm_ssh_user}"
}
output "dm_ssh_port" {
  value = "${var.dm_ssh_port}"
}
```

## Installing the driver

Download the [latest release](https://github.com/tintoy/docker-machine-driver-terraform/releases) and place the provider executable in the same directory as `docker-machine` executable (or somewhere on your `PATH`).

## Building the driver

If you'd rather run from source, run `make dev` and then `source ./use-dev-driver.sh`. You're good to go :)

See [CONTRIBUTING.md](CONTRIBUTING.md) for more detailed information about building / modifying the driver.
