// Package tui realize TUI intrface
package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	choices = []string{}
	menu    = []string{
		"0-Регистрация",
		"1-Аутентификация(клиент)",
		"2-Просмотреть",
		"3-Показать пароли",
		"4-Добавление",
		"5-Обновление",
		"6-Получить файл/папка",
		"7-Удаление",
		"8-Помощь",
		"9-Выход",
	}
)

type (
	model struct {
		cursor int
		choice string
	}
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("Выберите требуемый пункт меню:\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func MainMenu(t string) (int, error) {
	if len(t) < 10 {
		choices = menu[:2]
	} else {
		choices = menu[:]
	}
	p := tea.NewProgram(model{})

	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok && m.choice != "" {
		return m.cursor, nil

	}
	return -1, nil
}
