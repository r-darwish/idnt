//go:build windows

package providers

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestAppwiz_GetApplications(t *testing.T) {
	var cleaner AppWiz
	apps, err := cleaner.GetApplications()
	require.NoError(t, err)
	require.NotEmpty(t, apps)
}

func TestAppWiz_RemoveApplication(t *testing.T) {
	var cleaner AppWiz
	apps, err := cleaner.GetApplications()
	require.NoError(t, err)
	require.NotEmpty(t, apps)

	appToRemove := os.Getenv("IDNT_APP_TO_REMOVE")
	require.NotEmpty(t, appToRemove)

	for _, app := range apps {
		if app.Name == appToRemove {
			err := cleaner.RemoveApplication(&app)
			require.NoError(t, err)
		}
	}
}
