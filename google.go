package main

import (
	// "fmt"
	"os"

    "golang.org/x/net/context"
    // "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/compute/v1"

	"github.com/nanobox-io/nanobox-provider-golang"
)

var defaultCredentials = provider.Credentials{
	"access-json": os.Getenv("DEFAULT_ACCESS_JSON"),
	"project": os.Getenv("DEFAULT_PROJECT"),
}

type GoogleCompute struct {
}

	// Meta() Metadata
	// Catalog() ([]ServerOption, error)
	// Verify(Credentials) (bool, error)
	// AddKey(Credentials, KeyOrder) (Key, error)
	// ListKeys(Credentials) ([]Key, error)
	// ShowKey(Credentials, string) (Key, error)
	// DeleteKey(Credentials, string) error
	// AddServer(Credentials, ServerOrder) (Server, error)
	// ListServers(Credentials) ([]Server, error)
	// ShowServer(Credentials, string) (Server, error)
	// DeleteServer(Credentials, string) error
	// RebootServer(Credentials, string) error
	// RestartServer(Credentials, string) error


func (gc GoogleCompute) client(creds provider.Credentials) (*compute.Service, error) {

	token := creds["access-json"]
	config, err := google.JWTConfigFromJSON([]byte(token), compute.ComputeScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background())

	return compute.New(client)
}

func (gc GoogleCompute) defaultClient() (*compute.Service, error) {
	return gc.client(defaultCredentials)
}

// firewall: {"allowed":[{"IPProtocol":"all"}],
// "creationTimestamp":"2017-05-04T08:48:22.276-07:00",
// "description":"what",
// "id":"6392147894168122281",
// "kind":"compute#firewall",
// "name":"guy",
// "network":"https://www.googleapis.com/compute/v1/projects/keen-jigsaw-165218/global/networks/default",
// "selfLink":"https://www.googleapis.com/compute/v1/projects/keen-jigsaw-165218/global/firewalls/guy",
// "sourceRanges":["0.0.0.0/0"],
// "targetTags":["guy"]}

func (gc GoogleCompute) setupFirewall(creds provider.Credentials) {
	comp, err := gc.client(creds)
	if err != nil {
		return
	}

	// exit early if there is already a nanobox firewall
	_, err = comp.Firewalls.Get(creds["project"], "nanobox").Do()
	if err == nil {
		return
	}

	rule := &compute.Firewall{
		Name: "nanobox",
		Description: "firewall acceptions for nanobox servers",
		SourceRanges: []string{"0.0.0.0/0"},
		Allowed: []*compute.FirewallAllowed{
			&compute.FirewallAllowed{IPProtocol: "all"},
		},
		TargetTags: []string{"nanobox"},
	}

	comp.Firewalls.Insert(creds["project"], rule).Do()
	
}