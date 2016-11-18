package main

/*
 * Main program entry-point
 * ------------------------
 */

import (
	"fmt"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/drivers/plugin"
	"os"
	"path"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("%s %s\n\n", path.Base(os.Args[0]), DriverVersion)

		return
	}

	plugin.RegisterDriver(
		&Driver{BaseDriver: &drivers.BaseDriver{
			SSHUser: "root",
			SSHPort: 22,
		}},
	)
}
