package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pokedex "gokedex/internal/http_get"
	"log"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	List      list.Model
	TextInput textinput.Model
	Pokemon   pokedex.Pokemon
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// get a randomly selected pokemon
func initialModel() model {
	items := []list.Item{}
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()

	return model{
		TextInput: ti,
		List:      list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return m.TextInput.Cursor.BlinkCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			pokemon, err := pokedex.GetPokemon(m.TextInput.Value())
			if err != nil {
				log.Fatal(err)
			}
			m.Pokemon = pokemon

			render := []list.Item{
				item{"Name", m.Pokemon.Name},
				item{"Type(s)", pokedex.FormatTypes(m.Pokemon)},
				item{"Weight", fmt.Sprintln(m.Pokemon.Weight/10, "kg")},
				item{"Height", fmt.Sprintln(m.Pokemon.Height/10, "m")},
				item{"Abilities", pokedex.FormatAbilities(m.Pokemon)},
				item{"Stats", pokedex.FormatStats(m.Pokemon)},
				item{"Moves", pokedex.FormatMoves(m.Pokemon)},
				item{"Games", pokedex.FormatGames(m.Pokemon)},
			}
			m.List.SetItems(render)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	m.List, cmd = m.List.Update(msg)
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s", m.TextInput.View(), m.List.View())
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
