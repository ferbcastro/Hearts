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

	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "-c":
		player.InitPlayer(true)
	case "-e":
		player.InitPlayer(false)
	}

	for {
		if !player.IsGameActive() {
			println("Game ended!")
			break
		}

		if player.NoCardsLeft() {
			if player.IsCardDealer() {
				player.DealCards()
			} else {
				player.GetCards()
			}
		}

		player.Play()
		if player.IsRoundMaster() {
			player.WaitForAllCards()
			player.InformRoundLoser()
		} else {
			player.WaitForResult()
		}
	}
}
