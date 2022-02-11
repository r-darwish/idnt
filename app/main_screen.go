package app

import (
	"github.com/charmbracelet/bubbles/key"
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

type listKeyMap struct {
	execute    key.Binding
	selectItem key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		selectItem: key.NewBinding(
			key.WithKeys("tab", " ", "x"),
			key.WithHelp("tab", "select item"),
		),
		execute: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "uninstall selected apps"),
		),
	}
}

type AppItem struct {
	app    providers.Application
	marked bool
}

func (a AppItem) FilterValue() string {
	return a.app.Name
}

func (a AppItem) Title() string {
	title := a.app.Name
	if a.marked {
		title = "> " + title
	}
	return title
}

func (a AppItem) Description() string {
	return a.app.ExtraInfo["Description"]
}

type mainScreenModel struct {
	list list.Model
	done bool
	keys *listKeyMap
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
				allApps = append(allApps, &AppItem{app, false})
			}
		}

		sort.Slice(allApps, func(i, j int) bool {
			return strings.ToLower(allApps[i].(*AppItem).app.Name) > strings.ToLower(allApps[j].(*AppItem).app.Name)
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

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.selectItem):
			item := m.list.SelectedItem().(*AppItem)
			item.marked = !item.marked

			index := m.list.Index()
			if index < len(m.list.Items())-1 {
				m.list.Select(index + 1)
			}
		}
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
	listKeys := newListKeyMap()
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select applications to uninstall"
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.execute,
			listKeys.selectItem,
		}
	}

	return mainScreenModel{
		list: l,
		done: false,
		keys: listKeys,
	}
}
