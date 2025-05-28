package main

import (
	"Src/TokenRing"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		return
	}

	var socket TokenRing.SockDgram
	TokenRing.InitSocket(&socket, args[0], args[1])
	for {
		TokenRing.Send(&socket, []byte("Hello World"))
		time.Sleep(time.Second)
	}
}
