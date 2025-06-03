package main

import (
	"Src/Hearts"
	"log"
	"os"
)

const NUM_PLAYERS = 4

func usage() {
	log.Println("Use -c to create ring; -e to enter ring")
}

func main() {
	var player Hearts.Player
	var cardDeliver bool

	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "-c":
		player.InitPlayer(true)
		cardDeliver = true
	case "-e":
		player.InitPlayer(false)
		cardDeliver = false
	}

}
