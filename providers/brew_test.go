//go:build linux

package providers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBrew_GetApplications(t *testing.T) {
	brew := Brew{}
	applications, err := brew.GetApplications()
	require.NoError(t, err)
	require.NotEmpty(t, applications)
}
