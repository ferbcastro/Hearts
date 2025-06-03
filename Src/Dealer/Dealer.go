package main

import (
	"Src/TokenRing"
  "fmt"
)

func InitDealer() {

}

func main() {

	var client TokenRing.TokenRingClient

  ips := make([]string, 4)
  for i := 0; i < 4; i++ {
    fmt.Print("Enter ip: ")
    _, err := fmt.Scanln(&ips[i])
    if err != nil {
      fmt.Println("Error reading input:", err)
      continue
    }

    fmt.Println("You typed:", ips[i])
  }

  client.CreateRing(ips)
}
