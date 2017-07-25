package main

import (
	"github.com/nanobox-io/nanobox-provider-golang"
)

func (gc GoogleCompute) Meta() provider.Metadata {
	return provider.Metadata{
		ID: "gc",
		Name: "Google Compute (beta)",
		NickName: "instance",
		DefaultRegion: "us-west1-a",
		DefaultPlan: "standard",
		DefaultSize: "n1-standard-1",
		Rebootable: true,
		Renamable: false,
		SSHAuthMethod: "key",
		SSHKeyMethod: "object", // reference or object
		SSHUser: "ubuntu",
		ExternalInterface: "",
		InternalInterface: "ens4",
		CredentialFields: []provider.CredentialField{
			provider.CredentialField{
				Key: "access-json",
				Label: "Access Json",
			},
			provider.CredentialField{
				Key: "project",
				Label: "Project",
			},
		},
		Instructions: "https://developers.google.com/identity/protocols/OAuth2ServiceAccount",
		BootstrapScript: "https://s3.amazonaws.com/tools.nanobox.io/bootstrap/ubuntu.sh",
	}
}