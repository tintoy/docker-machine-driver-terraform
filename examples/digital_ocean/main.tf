/*
 * Configuration for Docker Machine (digitalocean)
 * -----------------------------------------------
 */

# Docker Machine variables
variable "dm_client_ip"           { }
variable "dm_machine_name"        { }
variable "dm_ssh_user"            { }
variable "dm_ssh_port"            { }
variable "dm_ssh_public_key_file" { }

# Additional variables
variable "os_image"           { }
variable "region"             { }
variable "size"               { }

# Digital Ocean
provider "digitalocean" {
  # Access token comes from DIGITALOCEAN_TOKEN environment variable
}

# SSH key for Docker host
resource "digitalocean_ssh_key" "docker_machine" {
    name		= "${var.dm_machine_name}@docker-machine"
	public_key	= "${file(var.dm_ssh_public_key_file)}"
}

# Docker host
resource "digitalocean_droplet" "docker_machine" {
  name        = "${var.dm_machine_name}"
  region      = "${var.region}"
  size        = "${var.size}"
  image       = "${var.os_image}"

  ssh_keys    = [ "${digitalocean_ssh_key.docker_machine.fingerprint}" ]
}

# Outputs for Docker Machine
output "dm_machine_ip" {
  value = "${digitalocean_droplet.docker_machine.ipv4_address}"
}
output "dm_ssh_user" {
  value = "${var.dm_ssh_user}"
}
output "dm_ssh_port" {
  value = "${var.dm_ssh_port}"
}
