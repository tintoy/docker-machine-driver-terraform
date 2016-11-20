/*
 * Configuration for Docker Machine (aws)
 * --------------------------------------
 */

# Docker Machine variables
variable "dm_client_ip"           { }
variable "dm_machine_name"        { }
variable "dm_ssh_user"            { }
variable "dm_ssh_port"            { }
variable "dm_ssh_public_key_file" { }

# Additional variables
variable "region"                 { }

# AWS
provider "aws" {
  region = "${var.region}"

  # Access key and secret key come from AWS_ACCESS_KEY_ID and AWS_SECRET_KEY environment variables
}

# Look up OS image (AMI)
data "aws_ami" "ubuntu" {
  owners        = ["099720109477"] # Canonical
  most_recent   = true

  filter {
    name    = "name"
    values  = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"]
  }

  filter {
    name    = "virtualization-type"
    values  = ["hvm"]
  }
}

# Docker host
resource "aws_instance" "docker_machine" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"

  key_name          = "${aws_key_pair.docker_machine.key_name}"
  security_groups   = [ "${aws_security_group.docker_machine.name}" ]

  tags {
    Name = "${var.dm_machine_name}"
  }
}

# Security group (allow ingress for SSH and Docker API)
resource "aws_security_group" "docker_machine" {
  name          = "dm_${replace(var.dm_machine_name, "-", "_")}"
  description   = "Docker Machine (allow SSH and Docker API)"

  # SSH (inbound)
  ingress {
      from_port     = 22
      to_port       = 22
      protocol      = "tcp"
      cidr_blocks   = ["${var.dm_client_ip}/0"]
  }

  # Docker (inbound)
  ingress {
      from_port     = 2376
      to_port       = 2376
      protocol      = "tcp"
      cidr_blocks   = ["${var.dm_client_ip}/0"]
  }

  # All traffic (outbound)
  egress {
      from_port     = 0
      to_port       = 0
      protocol      = "-1"
      cidr_blocks   = ["0.0.0.0/0"]
  }
}

# SSH key pair
resource "aws_key_pair" "docker_machine" {
  key_name    = "${var.dm_machine_name}@docker-machine"
  public_key  = "${file(var.dm_ssh_public_key_file)}"
}

# Outputs for Docker Machine
output "dm_machine_ip" {
  value = "${aws_instance.docker_machine.public_ip}"
}
output "dm_ssh_user" {
  value = "ubuntu" # The Ubuntu AMI requires you to log in as "ubuntu", not "root"
}
output "dm_ssh_port" {
  value = "${var.dm_ssh_port}"
}
