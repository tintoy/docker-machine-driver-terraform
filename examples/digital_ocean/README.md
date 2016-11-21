## digitalocean provider

This example deploys a droplet with a public IP address.

### Prerequisites

None.

### Additional variables

* `os_image` - An identifier representing the OS image from which the droplet will be created
* `region` - An identifier representing the region (datacenter) where the droplet will be deployed
* `size` - An identifier representing the size for the new droplet

### Using this example

First, update [additional-variables.json](additional-variables.json) with appropriate values.

```bash
export DIGITALOCEAN_TOKEN=my_token
cd examples/digital_ocean
docker-machine create --driver terraform \
	--terraform-config $PWD/main.tf \
	--terraform-variables-from $PWD/additional-variables.json \
	hello-digital-ocean
```

### Artefacts

This configuration produces:

* `digitalocean_droplet.docker_machine` - The Docker host
