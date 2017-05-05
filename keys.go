package main

import (
"fmt"
	"github.com/nanobox-io/nanobox-provider-golang"
)

// AddKey(Credentials, KeyOrder) (Key, error)
func (gc GoogleCompute) AddKey(creds provider.Credentials, keyOrder provider.KeyOrder) (provider.Key, error) {
	return provider.Key{}, fmt.Errorf("not yet implemented")
}

// ListKeys(Credentials) ([]Key, error)
func (gc GoogleCompute) ListKeys(creds provider.Credentials) ([]provider.Key, error) {
	return []provider.Key{}, fmt.Errorf("not yet implemented")
}

// ShowKey(Credentials, string) (provider.Key, error)
func (gc GoogleCompute) ShowKey(creds provider.Credentials, id string) (provider.Key, error) {
	return provider.Key{}, fmt.Errorf("not yet implemented")
}

// DeleteKey(Credentials, string) error
func (gc GoogleCompute) DeleteKey(creds provider.Credentials, id string) error {
	return fmt.Errorf("not yet implemented")
}
