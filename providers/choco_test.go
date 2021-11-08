package providers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChoco_GetApplications(t *testing.T) {
	choco := Choco{}
	applications, err := choco.GetApplications()

	require.NoError(t, err)
	require.NotEmpty(t, applications)

}

func TestFillChocoInfo(t *testing.T) {
	app := Application{Name: "7zip"}
	app.ExtraInfo = map[string]string{}
	err := fillChocoInfo(&app)

	require.NoError(t, err)
	require.Contains(t, app.ExtraInfo, "URL")
	require.Contains(t, app.ExtraInfo, "Description")
}
