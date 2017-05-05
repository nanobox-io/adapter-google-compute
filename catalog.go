package main

import (
	"strings"

    "google.golang.org/api/compute/v1"

	"github.com/nanobox-io/nanobox-provider-golang"
)

type zone struct {
	machineTypes map[string][]*compute.MachineType 
}

func (gc GoogleCompute) Catalog() ([]provider.ServerOption, error) {

	service, err := gc.defaultClient()
	if err != nil {
		return nil, err
	}

	zoneAggrigate, err := service.MachineTypes.AggregatedList("keen-jigsaw-165218").Do()
	if err != nil {
		return nil, err
	}

	serverOptions := []provider.ServerOption{}

	zonedata := map[string]zone{}

	// loop through all the zones
	for _, machine := range zoneAggrigate.Items {
		// loop through each machine type on in that zone
		for _, mt := range machine.MachineTypes {
			z, ok := zonedata[mt.Zone]
			// if it is a new zone
			if !ok {
				z = zone{machineTypes: map[string][]*compute.MachineType{}}
				zonedata[mt.Zone] = z
			}

			machineBucket := strings.Split(mt.Name, "-")[1]
			machineTypes, ok := z.machineTypes[machineBucket]
			if !ok {
				machineTypes =  []*compute.MachineType{}
				z.machineTypes[machineBucket] = machineTypes
			}

			machineTypes = append(machineTypes, mt)
			z.machineTypes[machineBucket] = machineTypes
		}
	}

	for zoneName, zone := range zonedata {
		serverOption := provider.ServerOption{ID: zoneName, Name: zoneName, Plans: []provider.ServerPlan{}}
		for bucket, mts := range zone.machineTypes {
			plan := provider.ServerPlan{ID: bucket, Name: bucket, Specs: []map[string]interface{}{}}
			for _, mt := range mts {
				plan.Specs = append(plan.Specs, map[string]interface{}{
					"id": mt.Name,
					"cpus": mt.GuestCpus,
					"ram": mt.MemoryMb / 1024.0,
					"disk": calculateDisk(int(mt.MemoryMb)),
					})
			}
			serverOption.Plans = append(serverOption.Plans, plan)
		}
		serverOptions = append(serverOptions, serverOption)
	}

	return serverOptions, nil
}

// use the amount of ram in each machine type to determine 
// disk size
func calculateDisk(ram int) int {
	gbs := ram / 1024
    switch {
    case gbs < 1:
    	return 20
    case gbs < 2:
    	return 30
    case gbs < 4:
    	return 40
    case gbs < 8:
    	return 60
    default:
    	return gbs * 10
    }
	return gbs * 10
}