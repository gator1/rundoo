// +build test

package registry

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/jarcoal/httpmock"
)

func TestServer(t *testing.T) {
	httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    expectedResponse := "Mocked response"
    httpmock.RegisterResponder("POST", "", httpmock.NewStringResponder(200, expectedResponse))


	registry := NewRegistry()
    r := ServiceConfig{Name: "TestService"}
	registry.RegisterService(r)

	err := registry.StartServices()
		

	require.NoError(t, err)

	err = registry.StopServices()
		

	require.NoError(t, err)

	r1 := ServiceConfig{Name: "TestService1"}
	err = registry.add(r1)
	require.NoError(t, err)
	
}
