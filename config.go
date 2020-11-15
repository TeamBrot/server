package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents parsed configuration data
type Config struct {
	Width   int
	Height  int
	Players int
}

func printUsageAndExit(status int) {
	fmt.Printf(`Usage: %s [OPTION] ...
Host a server for spe_ed 

  -h 	height of the board
  -p 	number of players (max: 63)
  -w	width of the board 
`, os.Args[0])
	os.Exit(status)
}

func parseInt(arg string, minValue int) int {
	i, err := strconv.Atoi(arg)
	if err != nil {
		printUsageAndExit(1)
	}
	if i < minValue {
		printUsageAndExit(1)
	}
	return i
}

// GetConfig parses the program arguments and returns the configuration
func GetConfig() Config {
	config := Config{40, 40, 2}

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-w":
			i++
			config.Width = parseInt(os.Args[i], 2)
			break
		case "-h":
			i++
			config.Height = parseInt(os.Args[i], 2)
			break
		case "-p":
			i++
			config.Players = parseInt(os.Args[i], 2)
			// We use a bitmask with the players' id, so the maximum id is 63
			if config.Players > 63 {
				printUsageAndExit(1)
			}
			break
		case "--help":
			printUsageAndExit(0)
		default:
			printUsageAndExit(1)
		}
	}
	return config
}
