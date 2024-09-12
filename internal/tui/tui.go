package tui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func PrintTable(results interface{}) {
	var rows []table.Row

	var columns []table.Column
	if objects, ok := results.([]types.Object); ok {
		columns = []table.Column{
			{Title: "Size", Width: 8},
			{Title: "Key", Width: 30},
			{Title: "Last Modified", Width: 30},
		}

		for _, object := range objects {
			rows = append(rows, []string{
				strconv.FormatInt(*object.Size, 10),
				*object.Key,
				object.LastModified.String(),
			})
		}
	} else if buckets, ok := results.([]types.Bucket); ok {
		columns = []table.Column{
			{Title: "Name", Width: 30},
			{Title: "Creation Date", Width: 30},
		}

		for _, bucket := range buckets {
			rows = append(rows, []string{
				*bucket.Name,
				bucket.CreationDate.String(),
			})
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(30),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
