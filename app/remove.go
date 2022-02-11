package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r-darwish/idnt/providers"
	"strings"
	"time"
)

type appForRemoval struct {
	app     providers.Application
	success bool
}

type removeModel struct {
	appsToRemove []*appForRemoval
	nextApp      int
}

func (r removeModel) Init() tea.Cmd {
	return tea.Batch(tea.ExitAltScreen, r.removeNextApp)
}

func (r removeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return r, tea.Quit
		}

	case appRemoved:
		r.appsToRemove[r.nextApp].success = msg.success
		r.nextApp += 1
		return r, r.removeNextApp

	case doneRemoving:
		return r, tea.Quit
	}

	return r, nil
}

func (r removeModel) View() string {
	builder := strings.Builder{}

	title := "Removing applications..."
	if r.nextApp >= len(r.appsToRemove) {
		title = "All Done!"
	}

	builder.WriteString(fmt.Sprintf("%s\n\n", titleStyle.Render(title)))

	for i, application := range r.appsToRemove {
		marker := " "
		if i < r.nextApp {
			if application.success {
				marker = "✅"
			} else {
				marker = "❌"
			}
		} else if i == r.nextApp {
			marker = "\U0001FA93"
		}
		builder.WriteString(fmt.Sprintf("%s %s\n\n", marker, application.app.Name))
	}

	return appStyle.Render(builder.String())
}

type appRemoved struct {
	success bool
}

type doneRemoving struct{}

func (r removeModel) removeNextApp() tea.Msg {
	if r.nextApp >= len(r.appsToRemove) {
		return doneRemoving{}
	}
	//app := r.appsToRemove[r.nextApp]
	time.Sleep(3 * time.Second)
	return appRemoved{success: true}
}

func newRemovalModel(appsForRemoval []providers.Application) tea.Model {
	modelData := make([]*appForRemoval, len(appsForRemoval))
	for i, application := range appsForRemoval {
		modelData[i] = &appForRemoval{app: application, success: false}
	}
	return removeModel{appsToRemove: modelData, nextApp: 0}
}