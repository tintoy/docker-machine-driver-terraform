## ddcloud provider (public IP address)

This configuration deploys a server and other resources required to expose it via a public IPv4 address (externally accessible only from the local client's IP).

An SSH public key is installed on the server, and password authentication is disabled.

Source:

* [Configuration](main.tf)
* [Additional Variables](additional-variables.json)

### Prerequisites

* An existing network domain
* An existing VLAN

### Additional variables

* `mcp_region` - The CloudControl region code that identifies the API end-point to use
* `networkdomain` - The name of the target network domain
* `datacenter` - The name of the datacenter where the network domain is hosted
* `vlan` - The name of the target VLAN within the network domain
* `ssh_bootstrap_password` - The temporary password used to bootstrap SSH authentication

### Using this example

First, update [additional-variables.json](additional-variables.json) with appropriate values.

```bash
cd examples/ddcloud/public_ipv4
export MCP_USERNAME=my_username
export MCP_PASSWORD=my_password
docker-machine create --driver terraform \
	--terraform-config $PWD/main.tf \
	--terraform-additional-variables $PWD/additional-variables.json \
	hello-ddcloud
```

### Artefacts

This configuration produces:

* `ddcloud_server.docker_machine` - The Docker host
* `ddcloud_nat.docker_machine` - The NAT rule used to route IPv4 traffic from a public IP address to the server's private IP address
* `ddcloud_firewall_rule.docker_machine_ssh4_in` - The firewall rule that allows incoming SSH connections from the client to the server's public IPv4 address  
Required by Docker Machine for provisioning; if you don't need it, you can disable this rule once provisioning is complete
* `ddcloud_firewall_rule.docker_machine_docker_in` - The firewall rule that allows incoming Docker API connections from the client to the server's public IPv4 address  
Required by Docker Machine for provisioning; if you don't need it, you can disable this rule once provisioning is complete
