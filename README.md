# GOkedex
*WIP*
<br>
This project was developed to learn the go programming language. 
<br>
I utilized a GET request to retrieve data from [this](https://pokeapi.co/) open-source API, process the information and display it using [BubbleTea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

---

## Building & Running 
With GoLang installed:
```
go build -o gokedex cmd/gokedex/main.go
```
In a linux or MacOs terminal:
```
./gokedex
```
## Possible Changes
1. The color of the text changes based on the type of the pokemon selected, for example serching pikachu will change the layout color to yellow, ponyta red and so on..
However this may cause problem when selecting **Dark** type pokemon.
2. In some fields, the data is hidden or truncated due to its length. It would be nice to open a popup showing the entire content of the field.
