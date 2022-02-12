package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/r-darwish/idnt/providers"
	"strings"
)

var (
	removalError = false
	normalStyle  = lipgloss.NewStyle().UnsetForeground()
	errorStyle   = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))
	doneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))
	removingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))
)

type appForRemoval struct {
	app     providers.Application
	success bool
}

type removeModel struct {
	appsToRemove []*appForRemoval
	nextApp      int
	spinner      spinner.Model
}

func (r removeModel) Init() tea.Cmd {
	return tea.Batch(tea.ExitAltScreen, r.removeNextApp, r.spinner.Tick)
}

func (r removeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return r, tea.Quit
		}

	case appRemoved:
		if !msg.success {
			removalError = true
		}
		r.appsToRemove[r.nextApp].success = msg.success
		r.nextApp += 1
		return r, r.removeNextApp

	case doneRemoving:
		r.spinner.Finish()
		return r, tea.Quit

	default:
		var cmd tea.Cmd
		r.spinner, cmd = r.spinner.Update(msg)
		return r, cmd
	}

	return r, nil
}

func (r removeModel) View() string {
	builder := strings.Builder{}

	title := "Removing applications..."
	done := r.nextApp >= len(r.appsToRemove)
	if done {
		title = "All Done!"
	}

	builder.WriteString(fmt.Sprintf("%s %s\n\n", r.spinner.View(), titleStyle.Render(title)))

	for i, application := range r.appsToRemove {
		marker := " "
		style := normalStyle

		if i < r.nextApp {
			if application.success {
				style = doneStyle
				marker = "-"
			} else {
				style = errorStyle
				marker = "!"
			}
		} else if i == r.nextApp {
			style = removingStyle
			marker = ">"
		}
		builder.WriteString(style.Render(fmt.Sprintf("%s %s", marker, application.app.Name)))
		builder.WriteString("\n\n")
	}

	if done && removalError {
		builder.WriteString(fmt.Sprintf("\n\n%s\n", errorStyle.Render("Failed removing some applications")))
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
	app := r.appsToRemove[r.nextApp]
	err := app.app.Provider.RemoveApplication(&app.app)
	return appRemoved{success: err == nil}
}

func newRemovalModel(appsForRemoval []providers.Application) tea.Model {
	modelData := make([]*appForRemoval, len(appsForRemoval))
	for i, application := range appsForRemoval {
		modelData[i] = &appForRemoval{app: application, success: false}
	}
	s := spinner.New()
	s.Spinner = spinner.Pulse
	return removeModel{
		appsToRemove: modelData,
		nextApp:      0,
		spinner:      s,
	}
}
