package main

import (
	"fmt"
	"strings"

    "google.golang.org/api/compute/v1"
	"github.com/nanobox-io/nanobox-provider-golang"
)

func (gc GoogleCompute) AddServer(creds provider.Credentials, order provider.ServerOrder) (provider.Server, error) {
	fmt.Printf("creds : %#v\n", creds)
	fmt.Printf("dreds : %#v\n", defaultCredentials)
	fmt.Printf("order : %#v\n", order)
	fmt.Println("sames:", creds["access-json"] == defaultCredentials["access-json"])

	// ensure firewall rule is created 
	gc.setupFirewall(creds)

	// take . out of names
	order.Name = strings.Replace(order.Name, ".", "-", -1)

	instance := &compute.Instance{
		Name: order.Name,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", order.Region, order.Size),
		Tags: &compute.Tags{Items: []string{"nanobox"}},
		Disks: []*compute.AttachedDisk{
			&compute.AttachedDisk{
				AutoDelete: true,
				Boot: true,
				DeviceName: order.Name,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/ubuntu-os-cloud/global/images/family/ubuntu-1604-lts",
					// DiskType: "pd-standard",
					DiskSizeGb: 10, // needs dynamic 
				},		
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			&compute.NetworkInterface{
				Network: fmt.Sprintf("projects/%s/global/networks/default", creds["project"]),
				AccessConfigs: []*compute.AccessConfig{
					&compute.AccessConfig{
						Name: "External NAT",
						Type: "ONE_TO_ONE_NAT",
					},
				},
			},
		},
	}

	if order.SSHKey != "" {
		val := fmt.Sprintf("ubuntu:%s", order.SSHKey)

		instance.Metadata = &compute.Metadata{
			Items: []*compute.MetadataItems{
				&compute.MetadataItems{
					Key: "ssh-keys",
					Value: &val,
				},
			},
		}
	}

	j, err := instance.MarshalJSON()
	fmt.Printf("instance: %s\nerr: %s\n", j, err)
	// get the machine type
	// create the disk
	// create server
	compute, err := gc.client(creds)
	if err != nil {
		fmt.Println("client:", err)
		return provider.Server{}, err
	}

	op, err := compute.Instances.Insert(creds["project"], order.Region, instance).Do()
	if err != nil {
		fmt.Println("insert:", err)
		return provider.Server{}, err
	}

	fmt.Printf("operation: %#v\n", op)

	server := provider.Server{
		ID: order.Name,
		Name: order.Name,
		Status: "created",
	}

	return server, nil
}