package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {
	commands := getCommands()
	cfg := config{
		cache: *pokecache.NewCache(1 * time.Minute),
		next:  "https://pokeapi.co/api/v2/location-area/",
		prev:  "",
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		sanitized := cleanInput(text)

		command, ok := commands[sanitized[0]]
		if !ok {
			fmt.Printf("Uknown command: %s\n", sanitized[0])
			continue
		}

		err := command.callback(&cfg)
		if err != nil {
			continue
		}

	}

}
