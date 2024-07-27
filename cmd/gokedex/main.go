package main

import (
	"fmt"
	pokedex "gokedex/internal/http_get"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	ResultText string
	Styles     *ComponentStyles
	Width      int
	Height     int
	List       list.Model
	TextInput  textinput.Model
	Pokemon    pokedex.Pokemon
	CanShow    bool
	Help       help.Model
	Keys       keyMap
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type keyMap struct {
	Up   key.Binding
	Down key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit},
	}
}

// get a randomly selected pokemon
func initialModel() model {
	items := []list.Item{}
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()

	m := model{
		TextInput:  ti,
		List:       list.New(items, list.NewDefaultDelegate(), 0, 0),
		CanShow:    false,
		Keys:       keys,
		Help:       help.New(),
		ResultText: "\n\n",
	}

	m.List.SetShowTitle(false)
	m.List.SetShowStatusBar(false)
	m.List.SetShowPagination(true)
	m.List.SetShowHelp(false)

	m.Styles = DefaultStyles()
	return m
}

func DefaultStyles() *ComponentStyles {
	s := new(ComponentStyles)
	return s
}

func (m model) Init() tea.Cmd {
	return m.TextInput.Cursor.BlinkCmd()
}

func getCorrectColor(el string) lipgloss.Color {
	switch el {
	case "normal":
		return lipgloss.Color("#FFF8F3")
	case "fire":
		return lipgloss.Color("#FF4C4C")
	case "water":
		return lipgloss.Color("#3FA2F6")
	case "electric":
		return lipgloss.Color("#FFDE4D")
	case "grass":
		return lipgloss.Color("#88D66C")
	case "ice":
		return lipgloss.Color("#96C9F4")
	case "fighting":
		return lipgloss.Color("#8C3061")
	case "poison":
		return lipgloss.Color("#4F1787")
	case "ground":
		return lipgloss.Color("#914F1E")
	case "flying":
		return lipgloss.Color("#36C2CE")
	case "psychic":
		return lipgloss.Color("#FF4191")
	case "bug":
		return lipgloss.Color("#508D4E")
	case "rock":
		return lipgloss.Color("#D6BD98")
	case "ghost":
		return lipgloss.Color("#180161")
	case "dragon":
		return lipgloss.Color("#03346E")
	case "dark":
		return lipgloss.Color("#201E43")
	case "steel":
		return lipgloss.Color("#EEEEEE")
	case "fairy":
		return lipgloss.Color("#FFAAAA")
	default:
		return lipgloss.Color("#FF4C4C")
	}

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			pokemon, err := pokedex.GetPokemon(strings.ToLower(m.TextInput.Value()))
			if err != nil {
				m.CanShow = true
			}
			m.Pokemon = pokemon
			m.CanShow = true
			render := []list.Item{
				item{"Name", strings.Title(m.Pokemon.Name)},
				item{"Type(s)", pokedex.FormatTypes(m.Pokemon)},
				item{"Weight", fmt.Sprintln(m.Pokemon.Weight/10, "kg")},
				item{"Height", fmt.Sprintln(m.Pokemon.Height/10, "m")},
				item{"Abilities", pokedex.FormatAbilities(m.Pokemon)},
				item{"Stats", pokedex.FormatStats(m.Pokemon)},
				item{"Moves", pokedex.FormatMoves(m.Pokemon)},
				item{"Games", pokedex.FormatGames(m.Pokemon)},
			}
			d := list.NewDefaultDelegate()
			var c lipgloss.Color
			if m.Pokemon.Name != "" {
				c = getCorrectColor(m.Pokemon.Types[0].Type.Name)
				m.ResultText = "\n\n"
			} else {
				m.CanShow = false
				m.ResultText = "\n\nPokemon not found\n\n"
			}
			d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
			d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy()

			m.List.SetDelegate(d)
			m.List.SetItems(render)

		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
		m.Width = msg.Width
		m.Height = msg.Height
	}

	m.List, cmd = m.List.Update(msg)
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

type ComponentStyles struct {
	List      lipgloss.Style
	ListColor lipgloss.Color
}

func (m model) View() string {
	helpView := m.Help.View(m.Keys)
	// Change colors
	if m.CanShow {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, m.TextInput.View(), m.Styles.List.Render(m.List.View()), helpView),
		)
	} else {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, m.TextInput.View(), m.ResultText, helpView),
		)
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
