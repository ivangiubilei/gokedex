package httpget

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TODO: aggiungere forze e debolezze
type Pokemon struct {
	Name   string  `json:"name"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`

	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`

	Games []struct {
		Version struct {
			Name string `json:"name"`
		} `json:"version"`
	} `json:"game_indices"`

	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
	} `json:"moves"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`

	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func (p Pokemon) formatTypes() string {
	ptype := "["
	for i, v := range p.Types {
		if str := strings.Replace(strings.Title(v.Type.Name), "-", " ", -1); i == len(p.Types)-1 {
			ptype += fmt.Sprintf("%s", str)
		} else {
			ptype += fmt.Sprintf("%s, ", str)
		}
	}
	return ptype + "]"
}

func (p Pokemon) formatStats() string {
	stats := "["
	for _, v := range p.Stats {
		str := strings.Replace(strings.Title(v.Stat.Name), "-", " ", -1)
		stats += fmt.Sprintf("\n\t%s: %d", str, v.BaseStat)
	}
	return stats + "\n]"
}

func (p Pokemon) formatMoves() string {
	moves := "["
	for i, v := range p.Moves {
		str := strings.Replace(strings.Title(v.Move.Name), "-", " ", -1)
		moves += fmt.Sprintf("\n\t(%d) %s", i+1, str)
	}
	return moves + "\n]"
}

func (p Pokemon) formatAbilities() string {
	abilities := "["
	for i, v := range p.Abilities {
		str := strings.Replace(strings.Title(v.Ability.Name), "-", " ", -1)
		if i == len(p.Abilities)-1 {
			abilities += str
		} else {
			abilities += str + ", "
		}
	}
	return abilities + "]"
}

func (p Pokemon) formatGames() string {
	games := "["
	for i, v := range p.Games {
		str := strings.Replace(strings.Title(v.Version.Name), "-", " ", -1)
		if i == len(p.Games)-1 {
			games += str
		} else {
			games += str + ", "
		}
	}
	return games + "]"
}

func (p Pokemon) String() string {
	return fmt.Sprintf("Name: %s\nTypes: %s\nHeight: %.1fm\nWeight: %.1fkg\nAbilities: %s\nBase stats: %s \nGames: %s\nMoves: %s", p.Name, p.formatTypes(), p.Height/10, p.Weight/10, p.formatAbilities(), p.formatStats(), p.formatGames(), p.formatMoves())
}

func GetPokemon(name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	client := http.Client{Timeout: time.Second * 3}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	req.Header.Set("User-Agent", "gokedex")

	res, getErr := client.Do(req)
	if getErr != nil {
		return Pokemon{}, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(io.Reader(res.Body))
	if readErr != nil {
		return Pokemon{}, readErr
	}

	pokemon := Pokemon{}

	jsonErr := json.Unmarshal(body, &pokemon)
	if jsonErr != nil {
		return Pokemon{}, errors.New("Pokemon not found")
	}

	return pokemon, nil
}
