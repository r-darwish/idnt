package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r-darwish/idnt/providers"
	"strings"
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
	return tea.ExitAltScreen
}

func (r removeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return r, nil
}

func (r removeModel) View() string {
	builder := strings.Builder{}

	for _, application := range r.appsToRemove {
		builder.WriteString(fmt.Sprintf("%s\n", application.app.Name))
	}

	return appStyle.Render(builder.String())
}

func (r removeModel) removeNextApp() tea.Msg {
	return nil
}

func newRemovalModel(appsForRemoval []providers.Application) tea.Model {
	modelData := make([]*appForRemoval, len(appsForRemoval))
	for i, application := range appsForRemoval {
		modelData[i] = &appForRemoval{app: application, success: false}
	}
	return removeModel{appsToRemove: modelData, nextApp: 0}
}
