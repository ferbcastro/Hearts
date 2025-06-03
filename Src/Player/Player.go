package main

import (
	"Src/Hearts"
	"Src/TokenRing"
	"fmt"
	"log"
	"os"
)

func usage() {
	log.Println("Use -c to create ring; -e to enter ring")
}

func main() {
	var client TokenRing.TokenRingClient
	var player Hearts.Player
	var ip string

	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "-c":
		ips := make([]string, 4)
		for i := 0; i < 4; i++ {
			fmt.Print("Enter ip: ")
			fmt.Scanln(&ips[i])
			fmt.Print("Enter order: ")
			fmt.Scanln()
		}
		client.CreateRing(ips)
	case "-e":
		fmt.Print("Enter ip: ")
		fmt.Scanln(&ip)
		client.EnterRing(ip)
	}

}
