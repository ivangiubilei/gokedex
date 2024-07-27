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
	MoreState  bool
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
	Up      key.Binding
	Down    key.Binding
	Quit    key.Binding
	Options key.Binding
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
	Options: key.NewBinding(
		key.WithKeys("options"),
		key.WithHelp("ctrl+g", "more options"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Options, k.Quit}
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
		MoreState:  false,
	}

	m.List.SetShowTitle(false)
	m.List.SetShowStatusBar(false)
	m.List.SetShowPagination(true)
	m.List.SetShowHelp(false)

	m.Styles = DefaultStyles()
	return m
}

func (i item) String() string {
	s := ""
	switch i.title {
	case "Stats":
		str := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(i.desc, "pecial", "."), " ", ""), ":", ": "), "[]")
		for i, el := range strings.Split(str, ",") {
			s += fmt.Sprintf("(%d) %s\n", i+1, el)
		}

	case "Type(s)":
		str := strings.Trim(i.desc, "[]")
		for i, el := range strings.Split(str, ", ") {
			s += fmt.Sprintf("(%d) %s\n", i+1, el)
		}
		// add data of the type

	case "Abilities":
		s = ""
		str := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(i.desc, " ", ""), ",", " "), "[]")
		for _, el := range strings.Split(str, " ") {
			s += fmt.Sprintf("%s\n", el)
		}

	case "Moves":
		s = ""
		count := strings.Count(i.desc, ", ")
		cols := 0
		str := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(i.desc, " ", ""), ",", " "), "[]")

		if count < 5 {
			cols = 1
		} else if count < 20 {
			cols = 2
		} else if count < 60 {
			cols = 3
		} else if count < 100 {
			cols = 4
		} else if count < 160 {
			cols = 5
		} else {
			cols = 6
		}

		for i, el := range strings.Split(str, " ") {

			if i%cols == 0 {
				s += "\n"
			}
			s += fmt.Sprintf("%s  ", el)
		}

	case "Games":
		s = ""
		str := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(i.desc, " ", ""), ",", " "), "[]")
		for i, el := range strings.Split(str, " ") {
			if i%2 == 0 {
				s += "\n"
			}
			s += fmt.Sprintf("%s  ", el)
		}
	case "Name":
		s += strings.Trim(i.desc, "[]")

	default:
		return fmt.Sprintln(i.desc)
	}
	return s
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
		return lipgloss.Color("#8C3061")
	case "dragon":
		return lipgloss.Color("#03346E")
	case "dark":
		return lipgloss.Color("#7C00FE")
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
		case "ctrl+g":
			if m.Pokemon.Name != "" {
				m.MoreState = !m.MoreState
			}
			m.View()
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

	if m.MoreState {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, fmt.Sprintln(m.List.SelectedItem())),
		)
	}

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
