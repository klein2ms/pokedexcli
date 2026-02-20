package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedexcli/internal/pokecache"
	"strings"
)

type config struct {
	cache pokecache.Cache
	next  string
	prev  string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 Pokemon locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 Pokemon locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays pokemon in the specified location",
			callback:    commandExplorer,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(cfg *config, args []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("	Usage:")

	for k, v := range getCommands() {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

type locationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type locationAreaResponse struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int           `json:"min_level"`
				MaxLevel        int           `json:"max_level"`
				ConditionValues []interface{} `json:"condition_values"`
				Chance          int           `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandMap(cfg *config, args []string) error {
	return printLocations(cfg, cfg.next)
}

func commandMapb(cfg *config, args []string) error {
	return printLocations(cfg, cfg.prev)
}

func printLocations(cfg *config, url string) error {

	var locations locationResponse

	results, ok := cfg.cache.Get(url)
	if ok {
		fmt.Println("Reading from cache...")
		err := json.Unmarshal(results, &locations)
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(res.Body)

		err = json.NewDecoder(res.Body).Decode(&locations)
		if err != nil {
			return err
		}

		bytes, err := json.Marshal(locations)
		if err != nil {
			return err
		}

		err = cfg.cache.Add(url, bytes)
		if err != nil {
			return err
		}
	}

	cfg.next = locations.Next
	cfg.prev = locations.Previous

	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandExplorer(cfg *config, args []string) error {

	location := args[0]

	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	var locationAreaResponse locationAreaResponse
	err = json.NewDecoder(res.Body).Decode(&locationAreaResponse)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationAreaResponse.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}
