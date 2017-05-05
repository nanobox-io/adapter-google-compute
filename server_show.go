package main

import (
	"fmt"
	
	"github.com/nanobox-io/nanobox-provider-golang"
)

// using the list function because we cant use the get functionality with out a zone provided
func (gc GoogleCompute) ShowServer(creds provider.Credentials, id string) (provider.Server, error) {
	servers, err := gc.ListServers(creds)
	if err != nil {
		return provider.Server{}, err
	}
	for _, server := range servers {
		if server.ID == id {
			return server, nil
		}
	}

	return provider.Server{}, fmt.Errorf("server not found")
}


// this way requires the zone to be provided and we didnt want to put the zone in the credentials
// func (gc GoogleCompute) ShowServer(creds provider.Credentials, id string) (provider.Server, error) {
// 	server := provider.Server{}

// 	compute, err := gc.client(creds)
// 	if err != nil {
// 		return server, err 
// 	}

// 	instance, err := compute.Instance.Get(creds["project"], creds["zone"], id).Do()
// 	if err != nil {
// 		return server, err
// 	}

// 	server = provider.Server{
// 	ID: fmt.Sprintf("%d", instance.Id),
// 	Status: instance.Status,
// 	Name: instance.Name,
// 	}

// 	// set the running status to something odin likes
// 	if server.Status == "RUNNING" {
// 		server.Status = "active"
// 	}

// 	// get the networking if it exists
// 	if len(instance.NetworkInterfaces) > 0 {
// 		server.InternalIP = instance.NetworkInterfaces[0].NetworkIP
// 		if len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
// 			server.ExternalIP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
// 		}
// 	}

// 	return server, nil
// }