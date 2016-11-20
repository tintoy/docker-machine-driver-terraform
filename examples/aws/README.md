## aws provider

This example deploys an EC2 instance and a security group to permit inbound SSH + Docker API.

Source:

* [Configuration](main.tf)
* [Additional Variables](additional-variables.json)

### Prerequisites

None.

### Additional variables

* `region` - The AWS region code that identifies the API end-point to use

### Using this example

First, update [additional-variables.json](additional-variables.json) with appropriate values.

```bash
cd examples/aws
export AWS_ACCESS_KEY_ID=my_access_key_id
export AWS_SECRET_KEY=my_secret_key
docker-machine create --driver terraform \
	--terraform-config $PWD/main.tf \
	--terraform-additional-variables $PWD/additional-variables.json \
	hello-aws
```

### Artefacts

This configuration produces:

* `aws_instance.docker_machine` - The Docker host
* `aws_security_group.docker_machine` - The network security group
