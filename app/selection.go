package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
	"github.com/r-darwish/idnt/providers"
	"sort"
	"strings"
)

type listKeyMap struct {
	execute    key.Binding
	selectItem key.Binding
	gotoUrl    key.Binding
	search     key.Binding
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
		gotoUrl: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "go to the application website"),
		),
		search: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "search the application in DuckDuckGo"),
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
	return fmt.Sprintf("%s (%s)", a.app.ExtraInfo["Description"], a.app.Provider.GetName())
}

type mainScreenModel struct {
	list list.Model
	done bool
	keys *listKeyMap
}

type loading struct{}

func (m mainScreenModel) Init() tea.Cmd {
	return tea.Batch(gatherApps, func() tea.Msg { return loading{} })
}

func gatherApps() tea.Msg {
	var providersList []providers.Provider

	providersList = append(providersList, providers.Powershell{})
	providersList = append(providersList, providers.GetOsSpecificProviders()...)
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
		return strings.ToLower(allApps[i].(*AppItem).app.Name) < strings.ToLower(allApps[j].(*AppItem).app.Name)
	})

	return allApps
}

func (m mainScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loading:
		m.list.Title = "Gathering installed applications..."
		cmd := m.list.StartSpinner()
		return m, cmd

	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case []list.Item:
		m.done = true
		setItems := m.list.SetItems(msg)
		m.list.Title = "Select applications to uninstall:"
		m.list.StopSpinner()
		return m, setItems

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.selectItem):
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				item := selectedItem.(*AppItem)
				item.marked = !item.marked

				index := m.list.Index()
				if index < len(m.list.Items())-1 {
					m.list.Select(index + 1)
				}
			}

		case key.Matches(msg, m.keys.gotoUrl):
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				item := selectedItem.(*AppItem)
				if url, ok := item.app.ExtraInfo["URL"]; ok {
					err := browser.OpenURL(url)
					if err != nil {
						cmd := m.list.NewStatusMessage(fmt.Sprintf("Error opening a browser: %v", err))
						return m, cmd
					}
				} else {
					cmd := m.list.NewStatusMessage("This application does not have a URL")
					return m, cmd
				}
			}

		case key.Matches(msg, m.keys.search):
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				item := selectedItem.(*AppItem)
				err := browser.OpenURL(fmt.Sprintf("https://duckduckgo.com/?q=%s", item.app.Name))
				if err != nil {
					cmd := m.list.NewStatusMessage(fmt.Sprintf("Error opening a browser: %v", err))
					return m, cmd
				}

			}

		case key.Matches(msg, m.keys.execute):
			var appsToRemove []providers.Application
			for _, item := range m.list.Items() {
				item := item.(*AppItem)
				if item.marked {
					appsToRemove = append(appsToRemove, item.app)
				}
			}

			if len(appsToRemove) == 0 {
				if item := m.list.SelectedItem(); item != nil {
					item := item.(*AppItem)
					appsToRemove = append(appsToRemove, item.app)
				}
			}

			if len(appsToRemove) > 0 {
				rm := newRemovalModel(appsToRemove)
				return rm, rm.Init()
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

func (m mainScreenModel) View() string {
	return appStyle.Render(m.list.View())
}

func newMainScreen() tea.Model {
	listKeys := newListKeyMap()
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Styles.Title = titleStyle
	l.SetSpinner(spinner.Pulse)
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
