package main

import (
  "Src/TokenRing"
  "fmt"
)  

func InitPlayer() {

}

func main() {
  var client TokenRing.TokenRingClient
  var ip string

  fmt.Print("Enter ip: ")
  _, err := fmt.Scanln(&ip)
  if err != nil {
      fmt.Println("Error reading input:", err)
  }


  client.EnterRing(ip)

	type a struct {
		nome string
	}

	var b a
	client.Recv(&b)

	fmt.Printf("%+v\n", b)
}
