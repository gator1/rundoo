package registry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()
    r := ServiceConfig{Name: "TestService"}
	registry.RegisterService(r)

	err := registry.StartServices()
		

	require.NoError(t, err)

	err = registry.StopServices()
		

	require.NoError(t, err)
	
}
