package main

import (
	"fmt"
	// "time"
	"testing"

	// "github.com/nanobox-io/nanobox-provider-golang"
)

func TestCatalog(t *testing.T) {
	gc :=  GoogleCompute{}
	catalog, err := gc.Catalog()
	if err != nil {
		t.Error("catalog failure", err)
	}
	if len(catalog) == 0 {
		t.Error("catalog is too short")
	}
	// fmt.Printf("%+v\n", catalog)
}

func TestClient(t *testing.T) {
	gc :=  GoogleCompute{}

	service, err := gc.defaultClient()
	if err != nil {
		t.Error("something", err)
	}

	gc.setupFirewall(defaultCredentials)
	firewalls, _ := service.Firewalls.List("keen-jigsaw-165218").Do()
	for _, firewall := range firewalls.Items {
		b, _ := firewall.MarshalJSON()
		fmt.Printf("firewall: %s\n\n", b)
	}

	// zones, err := service.Zones.List("keen-jigsaw-165218").Do()
	// for _, zone := range zones.Items {
	// 	fmt.Printf("zone: %#v\n", *zone)
	// }
	// fmt.Printf("%#v\nerr: %v", zones, err)

	// machines, err := service.MachineTypes.AggregatedList("keen-jigsaw-165218").Do()
	// num := len(machines.Items)
	// fmt.Printf("%#v\nerr: %v\n\n", machines, err)
	// fmt.Println(len(machines.Items))
	// for thing, machine := range machines.Items {
	// 	num += len(machine.MachineTypes)
	// 	for _, mach := range machine.MachineTypes {
	// 		fmt.Printf("thing: %#v, machine: %#v\n", thing, *mach)
	// 	}
	// }
	// fmt.Println(num)

	// instance := &compute.Instance{
	// 	Name: "test1",
	// 	MachineTypes: "zones/us-west1-a/machineTypes/n1-standard-1",

	// }
	// thing, err := service.Instances.Insert("keen-jigsaw-165218", "us-west1-a", instancee)

	// instances, err := service.Instances.List("keen-jigsaw-165218", "us-west1-a").Do()
	// fmt.Printf("%#v\nerr: %v", instances, err)
	// for _, instance := range instances.Items {
	// 	fmt.Printf("instance: %#v\n", *instance)
	// }

	// zones, err := service.Zones.List("keen-jigsaw-165218").Do()
	// for _, zone := range zones.Items {
	// 	fmt.Printf("zone: %#v\n", *zone)
	// }
	// fmt.Printf("%#v\nerr: %v", zones, err)
}

// func TestServerCreate(t *testing.T) {
// 	order := provider.ServerOrder{
// 		Name: "test1",
// 		Region: "us-west1-a",
// 		Size: "n1-standard-1",
// 		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6gnVVdKEOa03ia0f5k+Bwsbbq4viYoktnzUxSdLNF4HdiGCXR153c9YA7g4sqCJhOCa5YL6vZ6XPXAhNamNYT3ZzPdZEZMOhEpcFPH2ED8lVOtPRm6wMrN2W48A/500ebTvpXNOHIlrptYV6IzEU4NFoMbcK+Vr3BzVNnurXC2WLyjXgEqAaY7lbkIZhoru4Y3PG+KnE0rWsSHsHmaRAI2On1uFv1n2ySqmfvej4IFOt6zMQAtFcSqhTSY33JRT1hOWmAChqnApyBlcM8mgB4bwpm+p8SxskwSpU2JwWkFL2E9rcwlGqPOuAliUVOKZbNnS08nPYFTYLFd/2e0h+v",
// 	}

// 	gc := GoogleCompute{}

// 	server, err := gc.AddServer(defaultCredentials, order)
// 	if err != nil {
// 		t.Error("add server failure", err)
// 	}

// 	for i := 0; i < 100; i++ {
// 		<- time.After(time.Second)
// 		server, err = gc.ShowServer(defaultCredentials, "test1")
// 		if server.Status == "active" {
// 			return 
// 		}
// 	}
// 	t.Error("server never became active (100 seconds)")
// }

// func TestServerList(t *testing.T) {
// 	gc := GoogleCompute{}
// 	servers, err := gc.ListServers(defaultCredentials)
// 	if len(servers) == 0 || err != nil {
// 		t.Error("server list failed")
// 	}

// }

// func TestServerShow(t *testing.T) {
// 	gc := GoogleCompute{}
// 	server, err := gc.ShowServer(defaultCredentials, "test1")
// 	if err != nil {
// 		t.Error("unable to show server %s", err)
// 	}
// 	fmt.Printf("%#v\n", server)
// }

// func TestServerDelete(t *testing.T) {
// 	gc := GoogleCompute{}
// 	err := gc.DeleteServer(defaultCredentials, "test1")
// 	if err != nil {
// 		t.Error("unable to delete server %s", err)
// 	}
// }

