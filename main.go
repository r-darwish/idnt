package main

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pterm/pterm"
	"github.com/r-darwish/idnt/app"
	"github.com/r-darwish/idnt/providers"
	"sort"
	"strings"
)

func main() {
	app.Run()
	return

	spinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Collecting installed applications")
	var providersList []providers.Provider

	providersList = append(providersList, providers.Powershell{})
	providersList = append(providersList, providers.GetOsSpecificProviders()...)
	var allApps []providers.Application

	for _, provider := range providersList {
		providerApps, err := provider.GetApplications()
		if err != nil {
			pterm.Warning.Printfln("Error collecting %s: %s", provider.GetName(), err)
			continue
		}
		allApps = append(allApps, providerApps...)
	}

	sort.Slice(allApps, func(i, j int) bool {
		return strings.ToLower(allApps[i].Name) > strings.ToLower(allApps[j].Name)
	})

	if spinner != nil {
		_ = spinner.Stop()
	}

	selections, err := fuzzyfinder.FindMulti(
		allApps,
		func(i int) string {
			return allApps[i].Name
		},
		fuzzyfinder.WithPromptString("Select applications to remove > "),
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			app := allApps[i]
			result := fmt.Sprintf("%s\n\nProvider: %s", app.Name, app.Provider.GetName())
			for field, value := range app.ExtraInfo {
				result += fmt.Sprintf("\n%s: %s", field, value)
			}
			return result
		}),
	)
	if err != nil {
		return
	}

	p, _ := pterm.DefaultProgressbar.WithTotal(len(selections)).WithTitle("Removing Applications").WithRemoveWhenDone(true).Start()

	for _, selection := range selections {
		app := allApps[selection]
		name := app.Name
		p.UpdateTitle("Removing " + name)
		err := app.Provider.RemoveApplication(&app)
		if err != nil {
			pterm.Error.Printfln("Removing %s: %s", name, err)
		} else {
			pterm.Success.Printfln("Removed %s", name)
		}
		p.Increment()
	}
}
