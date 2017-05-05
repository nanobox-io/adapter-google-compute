package main

import (
	"github.com/nanobox-io/nanobox-provider-golang"
)

func (gc GoogleCompute) Verify(cred provider.Credentials) (bool, error) {
	return true, nil
}