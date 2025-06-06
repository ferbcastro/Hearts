package main

import (
	"Src/TokenRing"
	"fmt"
	"log"
	"os"
)

const NUM_PLAYERS = 4

func usage() {
	log.Println("Use -c to create ring; -e to enter ring")
}

func main() {
	var player TokenRing.Player

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
			fmt.Println("Game ended!")
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
			if player.IsThereAWinner() {
				player.AnounceWinner()
			}
		} else {
			player.WaitForResult()
		}

		player.PrintPoints()
	}
}
