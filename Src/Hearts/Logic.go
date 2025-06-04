package Hearts

import (
	"Src/TokenRing"
	"fmt"
	"log"
	"math/rand"
)

const MAX_CARDS_PER_ROUND = 13
const TOTAL_CARDS = 52
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

/* Message types */
const (
	CARDS = iota
	NEXT
	GAME_WINNER
	POINTS_QUERY
	ROUND_LOSER
)

type message struct {
	msgType byte
	cards   [MAX_CARDS_PER_ROUND]card
}

type Player struct {
	ringClient    TokenRing.TokenRingClient
	cards         [MAX_CARDS_PER_ROUND]card
	clockWiseIds  []byte
	cardsCount    byte
	positionInIds byte
	isCardDealer  bool
	isRoundMaster bool
	isGameActive  bool
	points        int
	msg           message
}

func incModN(n int, num int) int {
	return (num + 1) % n
}

func printCards(cards []card) {

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
		/* should call DealerCard() */
		player.isCardDealer = true
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
		/* should call GetCards() */
		player.isCardDealer = false
	}
	player.points = 0
	player.cardsCount = 0
}

/* Should be called once by card Dealer */
func (player *Player) DealCards() {
	var numbers []int
	for i := 0; i < TOTAL_CARDS; i++ {
		numbers = append(numbers, i)
	}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	var i int
	j := 0
	pos := incModN(NUM_PLAYERS, int(player.positionInIds))
	/* send cards to players */
	for i = 0; i < TOTAL_CARDS-MAX_CARDS_PER_ROUND; i++ {
		player.msg.cards[j].suit = numbers[i] / len(ranks)
		player.msg.cards[j].rank = numbers[i] % len(ranks)
		if j == MAX_CARDS_PER_ROUND-1 {
			j = 0
			player.msg.msgType = CARDS
			player.ringClient.Send(player.clockWiseIds[pos], player.msg)
			pos = incModN(NUM_PLAYERS, int(player.positionInIds))
		} else {
			j++
		}
	}
	/* set Dealer's cards */
	for i := 0; i < MAX_CARDS_PER_ROUND; i++ {
		player.cards[i].suit = numbers[i] / len(ranks)
		player.cards[i].rank = numbers[i] % len(ranks)
		if player.cards[i].suit == clubs && player.cards[i].rank == two {
			player.isRoundMaster = true
		}
	}
}

/* Should be called once by players (except Dealer)
 * Return true if whoever got cards will be first player */
func (player *Player) GetCards() {
	player.ringClient.Recv(&player.msg)
	if player.msg.msgType != CARDS {
		log.Printf("Message type not expected [%v] in [GetCards]\n", player.msg.msgType)
	}
	player.cards = player.msg.cards
}

/* Show cards to player and let they choose */
func (player *Player) Play() {
	if player.isRoundMaster == true {

	} else {
		if player.msg.msgType != NEXT {
			println("Message type not expected [%v] in [Play]\n", player.msg.msgType)
		}
	}
}

func (player *Player) InformRoundLoser() {

}

/* Dealer sends a message to the round winner */
func (player *Player) AnounceWinner() {

}

/* Players */
func (player *Player) WaitForResult() {

}

func (player *Player) IsRoundMaster() bool {
	return player.isRoundMaster
}

func (player *Player) IsGameActive() bool {
	return player.isGameActive
}

func (player *Player) IsCardDealer() bool {
	return player.isCardDealer
}

func (player *Player) NoCardsLeft() bool {
	return player.cardsCount == 0
}
