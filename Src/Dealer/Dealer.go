package main

import (
	"Src/TokenRing"
	"os"
)

func InitDealer() {

}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		return
	}

	var socket TokenRing.SockDgram
}
