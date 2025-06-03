package Hearts

import (
	"Src/TokenRing"
	"fmt"
)

const MAX_CARDS = 13
const NUM_PLAYERS = 4

const (
	clubs = iota
	hearts
	spades
	diamonds
)

var suits = []string{"♣", "♥", "♠", "♦"}

const (
	two = iota
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	jack
	queen
	king
	ace
)

var ranks = []string{
	"𝔱𝔴𝔬",
	"𝔱𝔥𝔯𝔢𝔢",
	"𝔣𝔬𝔲𝔯",
	"𝔣𝔦𝔳𝔢",
	"𝔰𝔦𝔵",
	"𝔰𝔢𝔳𝔢𝔫",
	"𝔢𝔦𝔤𝔥𝔱",
	"𝔫𝔦𝔫𝔢",
	"𝔱𝔢𝔫",
	"𝔧𝔞𝔠𝔨",
	"𝔮𝔲𝔢𝔢𝔫",
	"𝔨𝔦𝔫𝔤",
	"𝔞𝔠𝔢",
}

type card struct {
	suit int
	rank int
}

type message struct {
	msgT  int
	cards [MAX_CARDS]card
}

type Player struct {
	ringClient    TokenRing.TokenRingClient
	cards         []byte
	clockWiseIds  []byte
	positionInIds byte
}

func printCard(rank int, suit int) {

}

func (player *Player) InitPlayer(isRingCreator bool) {
	var myId byte
	var ids []byte
	var ip string

	if isRingCreator == true {
		player.positionInIds = 0
		/* read ips */
		ips := make([]string, NUM_PLAYERS)
		fmt.Print("Enter your ip: ")
		fmt.Scanln(&ips[player.positionInIds])
		fmt.Println("Now enter other players' ip in clockwise order")
		for i := 1; i < NUM_PLAYERS; i++ {
			fmt.Print("Enter ip: ")
			fmt.Scanln(&ips[i])
		}
		/* creat ring */
		ids = player.ringClient.CreateRing(ips)
		player.clockWiseIds = ids
		/* broadcast ids */
		for i := 1; i < NUM_PLAYERS; i++ {
			player.ringClient.Send(ids[i], ids)
		}
	} else {
		/* enter ring */
		fmt.Print("Enter your ip: ")
		fmt.Scanln(&ip)
		myId = player.ringClient.EnterRing(ip)
		/* read ids */
		player.ringClient.Recv(&player.clockWiseIds)
		/* discover position */
		for i := byte(1); i < NUM_PLAYERS; i++ {
			if player.clockWiseIds[i] == myId {
				player.positionInIds = i
				break
			}
		}
	}
}

/* Should be called once by Dealer */
func (player *Player) DeliverCards() {

}

/* Should be called once by players (except Dealer) */
func (player *Player) GetCards() {

}

/* Show cards to player and let thy choose */
func (player *Player) Play() {

}

/* Dealer sends a message to the round winner */
func (player *Player) AnounceWinner() {

}

/* Players  */
func (player *Player) WaitForResult() {

}
