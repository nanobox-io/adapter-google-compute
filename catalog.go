package main

import (
	"strings"
	"fmt"
	"io"
	"encoding/json"
	"net/http"

    "google.golang.org/api/compute/v1"

	"github.com/nanobox-io/nanobox-provider-golang"
)

type zone struct {
	machineTypes map[string][]*compute.MachineType 
}

// a cash of the json pricing
var jsonCache map[string]map[string]interface{}

func (gc GoogleCompute) Catalog() ([]provider.ServerOption, error) {

	service, err := gc.defaultClient()
	if err != nil {
		return nil, err
	}

	zoneAggrigate, err := service.MachineTypes.AggregatedList("keen-jigsaw-165218").Do()
	if err != nil {
		fmt.Println(err)
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
				// get the pricing
				hourly, monthly := getPrice(zoneName, mt.Name)
				// create a new pricing spec
				newSpec := map[string]interface{}{
					"id": mt.Name,
					"cpu": mt.GuestCpus,
					"ram": mt.MemoryMb,
					"disk": calculateDisk(int(mt.MemoryMb)),
					"transfer": "unlimited",
					"dollars_per_hr": hourly,
					"dollars_per_mo": monthly,
				}

				plan.Specs = append(plan.Specs, newSpec)
			}
			serverOption.Plans = append(serverOption.Plans, plan)
		}
		serverOptions = append(serverOptions, serverOption)
	}

	// clean up the jsoncache for pricing
	defer func() {
		jsonCache = nil	
	}()
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

// given the api for pricing this should pull the data and get the
// price of the zone based on that data
func getPrice(zone, id string) (interface{}, interface{}) {
	pricing := priceJson()
	
	idName := fmt.Sprintf("CP-COMPUTEENGINE-VMIMAGE-%s", strings.ToUpper(id))
	price, ok := pricing[idName]
	if ok {
		zoneName := getZone(price, zone)
		hourly, ok := price[zoneName]
		if ok {
			houlyNumber, ok := hourly.(float64)
			if ok {
				return houlyNumber, houlyNumber * 720			
			}
		}
	}

	// if we cant find the price it is safer to say we dont know the price
	// then set it to 0.0 or something else
	return "unknown", "unknown"
}

// get the zone name. this may require us to strip the end character off the full zone until we get a match
func getZone(pricing map[string]interface{}, zoneFull string) string {
	shortZone := zoneFull
	for {
		// if we have nothing left in the short zone 
		// we can return the most generic string so asia-northeast1-c -> asia
		if shortZone == "" {
			return strings.Split(zoneFull, "-")[0]
		}

		// check to see if the short zone currently matches any pricing
		_, ok := pricing[shortZone]
		if ok {
			return shortZone
		}

		// if no match is found strip the last character and try again
		shortZone = shortZone[0:len(shortZone)-1]

	}
}


func priceJson() map[string]map[string]interface{} {
	if jsonCache != nil {
		return jsonCache
	}

	// pull the json from https://cloudpricingcalculator.appspot.com/static/data/pricelist.json
	resp, err := http.Get("https://cloudpricingcalculator.appspot.com/static/data/pricelist.json")
	if err != nil {
		return jsonCache
	}
	decode := json.NewDecoder(resp.Body)

	// create the json cache value
	jsonCache = map[string]map[string]interface{}{

	}

	for {
		token, err := decode.Token()
		// when we reach the end we are done
		if err == io.EOF {
			break
		}
		// an un expected error
		if err != nil {
			return jsonCache
		}

		// check to see if the value is a string
		str, ok := token.(string)
		// if it is a string and it is one of our compute pricing hashes
		// we need to add it to the json cache
		if ok && strings.Contains(str, "COMPUTEENGINE") {
			pricing := map[string]interface{}{}
			err := decode.Decode(&pricing)
			// if we are able to decode the pricing
			// add it
			if err == nil {
				jsonCache[str] = pricing
			}
		}
	}

	return jsonCache

}