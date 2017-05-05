package main

import (
	"fmt"

    "google.golang.org/api/compute/v1"
	"github.com/nanobox-io/nanobox-provider-golang"
)

func (gc GoogleCompute) ListServers(creds provider.Credentials) ([]provider.Server, error) {
	compute, err := gc.client(creds)
	if err != nil {
		return nil, err
	}

	instanceList, err := compute.Instances.AggregatedList(creds["project"]).Do()
	if err != nil {
		return nil, err
	}

	servers := []provider.Server{}
	for _, zoneInstanceList := range instanceList.Items {
		for _, instance := range zoneInstanceList.Instances {
			b, _ := instance.MarshalJSON()
			fmt.Printf("Instance: %s\n\n", b)

			server := provider.Server{
				ID: instance.Name,
				Status: instance.Status,
				Name: instance.Name,
				}

			// set the running status to something odin likes
			if server.Status == "RUNNING" {
				server.Status = "active"
			}

			// get the networking if it exists
			if len(instance.NetworkInterfaces) > 0 {
				server.InternalIP = instance.NetworkInterfaces[0].NetworkIP
				if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
					server.ExternalIP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
				}
			}

			servers = append(servers, server)
		}
		
	}

	return servers, nil
}

func (gc GoogleCompute) listInstances(creds provider.Credentials) ([]*compute.Instance, error) {
	service, err := gc.client(creds)
	if err != nil {
		return nil, err
	}

	instanceList, err := service.Instances.AggregatedList(creds["project"]).Do()
	if err != nil {
		return nil, err
	}

	instances := []*compute.Instance{}
	for _, zoneInstanceList := range instanceList.Items {
		for _, instance := range zoneInstanceList.Instances {
			instances = append(instances, instance)
		}
	}
	
	return instances, nil
}