package providers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPowershell_GetApplications(t *testing.T) {
	choco := Powershell{}
	applications, err := choco.GetApplications()

	require.NoError(t, err)
	require.NotEmpty(t, applications)

}
