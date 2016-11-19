/*
 * Configuration for Docker Machine (ddcloud, private IPv4 address)
 * ----------------------------------------------------------------
 */

# Docker Machine variables
variable "dm_client_ip"           { }
variable "dm_machine_name"        { }
variable "dm_ssh_user"            { }
variable "dm_ssh_port"            { }
variable "dm_ssh_public_key_file" { }

# Additional variables
variable "mcp_region"             { }
variable "networkdomain"          { }
variable "datacenter"             { }
variable "vlan"                   { }
variable "ssh_bootstrap_password" { }

# CloudControl
provider "ddcloud" {
  region    = "${var.region}"

  # Username and password come from MCP_USERNAME and MCP_PASSWORD environment variables
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

      host      = "${ddcloud_server.docker_machine.primary_adapter_ipv4}"
    }
  }
}

# Outputs for Docker Machine
output "dm_machine_ip" {
  value = "${ddcloud_server.docker_machine.primary_adapter_ipv4}"
}
output "dm_ssh_user" {
  value = "${var.dm_ssh_user}"
}
output "dm_ssh_port" {
  value = "${var.dm_ssh_port}"
}
