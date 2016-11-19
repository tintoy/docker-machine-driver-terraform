## ddcloud provider (private IP address)

This example deploys a server with a private IP address only; it needs to be run either while connected to the CloudControl VPN or from a server inside the target network.

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
cd examples/ddcloud/private_ipv4
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
