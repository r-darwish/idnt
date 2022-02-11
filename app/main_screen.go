package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/r-darwish/idnt/providers"
	"sort"
	"strings"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type AppItem providers.Application

func (a AppItem) FilterValue() string {
	return a.Name
}

func (a AppItem) Title() string {
	return a.Name
}

func (a AppItem) Description() string {
	return a.ExtraInfo["Description"]
}

type mainScreenModel struct {
	list list.Model
	done bool
}

func (m mainScreenModel) Init() tea.Cmd {
	return func() tea.Msg {
		var providersList []providers.Provider

		providersList = append(providersList, providers.Powershell{})
		// providersList = append(providersList, providers.GetOsSpecificProviders()...)
		var allApps []list.Item

		for _, provider := range providersList {
			providerApps, err := provider.GetApplications()
			if err != nil {
				continue
			}
			for _, app := range providerApps {
				allApps = append(allApps, AppItem(app))
			}
		}

		sort.Slice(allApps, func(i, j int) bool {
			return strings.ToLower(allApps[i].(AppItem).Name) > strings.ToLower(allApps[j].(AppItem).Name)
		})

		return allApps
	}
}

func (m mainScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case []list.Item:
		m.done = true
		cmd := m.list.SetItems(msg)
		return m, cmd
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

func (m mainScreenModel) View() string {
	if !m.done {
		return appStyle.Render("Loading...")
	} else {
		return appStyle.Render(m.list.View())
	}
}

func newMainScreen() tea.Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select applications to uninstall"
	return mainScreenModel{
		list: l,
		done: false,
	}
}
