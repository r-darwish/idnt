package main

import (
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pterm/pterm"
	"github.com/r-darwish/idnt/providers"
	"sort"
	"strings"
)

func main() {
	providersList := []providers.Provider{&providers.AppWiz{}}
	var allApps []providers.Application

	for _, provider := range providersList {
		providerApps, err := provider.GetApplications()
		if err != nil {
			continue
		}
		allApps = append(allApps, providerApps...)
	}

	sort.Slice(allApps, func(i, j int) bool {
		return strings.ToLower(allApps[i].Name) > strings.ToLower(allApps[j].Name)
	})

	selections, err := fuzzyfinder.FindMulti(
		allApps,
		func(i int) string {
			return allApps[i].Name
		},
	)
	if err != nil {
		return
	}

	p, _ := pterm.DefaultProgressbar.WithTotal(len(selections)).WithTitle("Removing Applications").Start()

	for _, selection := range selections {
		app := allApps[selection]
		name := app.Name
		p.UpdateTitle("Removing " + name)
		err := app.Provider.RemoveApplication(&app)
		p.Increment()
		if err != nil {
			pterm.Error.Printfln("Error removing %s: %s", name, err)
		} else {
			pterm.Success.Printfln("Removed %s", name)
		}
	}
}
