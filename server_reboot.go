package main

import (
	"strings"
    "google.golang.org/api/compute/v1"
	"github.com/nanobox-io/nanobox-provider-golang"
)

func (gc GoogleCompute) RebootServer(creds provider.Credentials, id string) error {
	instances, err := gc.listInstances(creds)
		if err != nil {
		return err
	}

	// select my instance
	var instance *compute.Instance
	for _, i := range instances {
		if i.Name == id {
			instance = i
			break
		}
	}

	// no instance found
	if instance == nil {
		return nil
	}

	// remove the found instance
	compute, err := gc.client(creds)
	if err != nil {
		return err
	}

	// Pull out the zone short name from the zone url
	// fmt.Printf("%#v\n\n", *instance)
	zoneParts := strings.Split(instance.Zone, "/zones/")
	zone := ""
	if len(zoneParts) == 2 {
		zone = zoneParts[1]
	}
	_, err = compute.Instances.Reset(creds["project"], zone, id).Do()
	if err != nil {
		return err
	}

	return nil
}