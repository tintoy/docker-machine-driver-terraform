/*
 * Configuration for Docker Machine (ddcloud, public IPv4 address)
 * ---------------------------------------------------------------
 */

# Docker Machine variables
variable "dm_client_ip"           { }
variable "dm_machine_name"        { }
variable "dm_ssh_user"            { }
variable "dm_ssh_port"            { }
variable "dm_ssh_public_key_file" { }

# Additional variables
variable "region"                 { }
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

# Public IPv4 address and firewall rules
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
