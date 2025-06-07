package main

import (
	"Src/Hearts"
	"fmt"
	"os"
)

const NUM_PLAYERS = 4

func usage() {
	fmt.Println("Use -c to create ring; -e to enter ring")
	fmt.Println("For ring creator: enter your ip followed by other players' ip",
		"in clockwise order")
	fmt.Println("For ring user: enter your ip")
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
	case "-h":
		usage()
	}

	for {
		if !player.IsGameActive() {
			fmt.Println("Game ended!")
			break
		}

		fmt.Printf("New round!\n\n")

		if player.NoCardsLeft() {
			player.ResetHeartsBroken()
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
			for {
				ret := player.WaitForResult()
				if ret == Hearts.ALL_RESULTS_GOT {
					break
				}
			}
		}

		player.PrintPoints()
	}
}
