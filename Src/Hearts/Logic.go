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
	CLUBS = iota
	HEARTS
	SPADES
	DIAMONDS
)

var suits = []string{"♣", "♥", "♠", "♦"}

const (
	TWO = iota
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	JACK
	QUEEN
	KING
	ACE
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
	suit uint8
	rank uint8
}

type deck []card

/* Message types */
const (
	CARDS = iota
	NEXT
	GAME_WINNER
	POINTS_QUERY
	ROUND_LOSER
	HEARTS_BROKEN
	BEGIN_GAME
)

type message struct {
	msgType        uint8
	cards          [MAX_CARDS_PER_ROUND]card
	numPlayedCards uint8
}

type Player struct {
	ringClient     TokenRing.TokenRingClient
	cards          [MAX_CARDS_PER_ROUND]card
	clockWiseIds   []uint8
	cardsCount     uint8
	positionInIds  uint8
	points         int
	msg            message
	isRoundMaster  bool
	isCardDealer   bool
	isGameActive   bool
	isHeartsBroken bool
}

func (c card) isCardEqual(rank, suit byte) bool {
	return (c.rank == rank && c.suit == suit)
}

func printCards(cards []card) {

}

func (m message) printUnexpectedType() {
	println("Message type not expected [%v] in [Play]\n", m.msgType)
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
	player.isHeartsBroken = false
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
	var j int
	var idNextRoundHead uint8
	pos := player.positionInIds
	for i = 0; i < TOTAL_CARDS; i++ {
		if (i%MAX_CARDS_PER_ROUND == 0) && (i != 0) {
			if i == MAX_CARDS_PER_ROUND { /* dealer cards */
				player.cards = player.msg.cards
			} else { /* other players' cards */
				player.msg.msgType = CARDS
				player.ringClient.Send(player.clockWiseIds[pos], player.msg)
			}
			pos = (pos + 1) % NUM_PLAYERS
		}

		j = i / MAX_CARDS_PER_ROUND
		player.msg.cards[j].suit = uint8(numbers[i] / len(ranks))
		player.msg.cards[j].rank = uint8(numbers[i] % len(ranks))
		if player.msg.cards[j].isCardEqual(TWO, CLUBS) {
			idNextRoundHead = player.clockWiseIds[pos]
		}
	}

	if idNextRoundHead == player.clockWiseIds[player.positionInIds] {
		player.isRoundMaster = true
	} else {
		player.msg.msgType = BEGIN_GAME
		player.ringClient.Send(idNextRoundHead, player.msg)
	}
}

/* Should be called once by players (except Dealer) */
func (player *Player) GetCards() {
	player.ringClient.Recv(&player.msg)
	if player.msg.msgType != CARDS {
		log.Printf("Message type not expected [%v] in [GetCards]\n", player.msg.msgType)
	}
	player.cards = player.msg.cards
	for i := 0; i < MAX_CARDS_PER_ROUND; i++ {
		if player.cards[i].isCardEqual(TWO, CLUBS) {
			player.isRoundMaster = true
		}
		player.ringClient.Recv(&player.msg)
		if player.msg.msgType != BEGIN_GAME {
			player.msg.printUnexpectedType()
		}
	}
}

/* Show cards to player and let they choose */
func (player *Player) Play() {
	printCards(player.cards[:])
	if player.isRoundMaster == true {
		fmt.Printf("You begin round!", player.positionInIds)
		if player.isHeartsBroken {

		}
	} else {
		if player.msg.msgType != NEXT {
			player.msg.printUnexpectedType()
		}
	}
}

/* Should be called by round master */
func (player *Player) InformRoundLoser() {

}

/* Should be called by round master */
func (player *Player) WaitForAllCards() {

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
