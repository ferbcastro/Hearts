package Hearts

import (
	"Src/TokenRing"
	"fmt"
	"math/rand"
)

const MAX_CARDS_PER_ROUND = 13
const TOTAL_CARDS = 52
const NUM_PLAYERS = 4
const MAX_POINTS = 50

const (
	DIAMONDS = iota
	SPADES
	HEARTS
	CLUBS
)

var suits = []string{"♦", "♠", "♥", "♣"}

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

/* Message types */
const (
	CARDS = iota
	NEXT
	GAME_WINNER
	POINTS_QUERY
	ROUND_LOSER
	CONTINUE_GAME
	HEARTS_BROKEN
	BEGIN_GAME
)

type message struct {
	cards          [MAX_CARDS_PER_ROUND]card
	msgType        uint8
	numPlayedCards uint8
	earnedPoints   uint8
}

type deck struct {
	cards     [MAX_CARDS_PER_ROUND]card
	wasUsed   [MAX_CARDS_PER_ROUND]bool
	cardsLeft uint8
}

type Player struct {
	ringClient     TokenRing.TokenRingClient
	clockWiseIds   []uint8
	positionInIds  uint8
	deck           deck
	msg            message
	points         int
	isRoundMaster  bool
	isCardDealer   bool
	isGameActive   bool
	isHeartsBroken bool
}

func (player *Player) InitPlayer(isRingCreator bool) {
	var myId byte
	var ids []byte
	var ip string

	if isRingCreator {
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
		player.positionInIds = 0
		/* broadcast ids */
		player.sendForAll(ids)
		/* should call DealerCard() */
		player.isCardDealer = true
	} else {
		/* enter ring */
		fmt.Print("Enter your ip: ")
		fmt.Scanln(&ip)
		myId = byte(player.ringClient.EnterRing(ip))
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
	player.deck.cardsLeft = 0
	player.isHeartsBroken = false
}

/* Should be called by card Dealer */
func (player *Player) DealCards() {
	var numbers []int
	for i := range TOTAL_CARDS {
		numbers = append(numbers, i)
	}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	fmt.Println("Cards shuffled! Distributing...")

	var j int
	var idNextRoundHead uint8
	pos := player.positionInIds
	for i := range TOTAL_CARDS {
		if (i%MAX_CARDS_PER_ROUND == 0) && (i != 0) {
			if i == MAX_CARDS_PER_ROUND { /* dealer cards */
				player.deck.initDeck(player.msg.cards)
			} else { /* other players' cards */
				player.msg.msgType = CARDS
				player.ringClient.Send(player.clockWiseIds[pos], player.msg)
			}
			pos = (pos + 1) % NUM_PLAYERS
		}

		j = i % MAX_CARDS_PER_ROUND
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

	fmt.Println("Players all set!")
}

/* Should be called by players (except Dealer) */
func (player *Player) GetCards() {
	fmt.Println("Getting cards...")
	player.ringClient.Recv(&player.msg)
	if player.msg.msgType != CARDS {
		player.msg.printUnexpectedType()
	}
	fmt.Println("Got cards!")
	player.deck.initDeck(player.msg.cards)
	for i := range player.deck.cards {
		if player.deck.cards[i].isCardEqual(TWO, CLUBS) {
			player.isRoundMaster = true
			player.ringClient.Recv(&player.msg)
			if player.msg.msgType != BEGIN_GAME {
				player.msg.printUnexpectedType()
			}
			break
		}
	}
}

/* Show cards to player and let they choose */
func (player *Player) Play() {
	var selected int
	var cardIt int
	switch player.isRoundMaster {
	case true:
		fmt.Println("You begin round!")
		player.deck.printDeck()
		for {
			fmt.Print("Choose a card from your deck: ")
			fmt.Scanln(&selected)
			cardIt = player.deck.getCardFromDeck(selected)
			if cardIt == -1 {
				fmt.Println("Invalid card!")
				continue
			}
			if player.deck.cards[cardIt].isSuitEqual(HEARTS) && !player.isHeartsBroken {
				fmt.Println("Invalid card!")
				continue
			}
			break
		}
	case false:
		player.ringClient.Recv(&player.msg)
		if player.msg.msgType == HEARTS_BROKEN {
			player.isHeartsBroken = true
			player.ringClient.Recv(&player.msg)
		}
		if player.msg.msgType != NEXT {
			player.msg.printUnexpectedType()
		}
		masterPos := (player.positionInIds + (NUM_PLAYERS - player.msg.numPlayedCards)) % NUM_PLAYERS
		masterSuit := player.msg.cards[masterPos].suit
		hasMasterSuit := false
		for i := range player.deck.cards {
			if !player.deck.wasUsed[i] && player.deck.cards[i].isSuitEqual(masterSuit) {
				hasMasterSuit = true
			}
		}
		fmt.Println("Your turn!")
		player.deck.printDeck()
		for {
			fmt.Print("Choose a card from your deck: ")
			fmt.Scanln(&selected)
			cardIt = player.deck.getCardFromDeck(selected)
			if cardIt == -1 {
				fmt.Println("Invalid card!")
				continue
			}
			if hasMasterSuit && !player.deck.cards[cardIt].isSuitEqual(masterSuit) {
				fmt.Println("Invalid card!")
				continue
			}
			if !player.isHeartsBroken && player.deck.cards[cardIt].isSuitEqual(HEARTS) {
				player.msg.msgType = HEARTS_BROKEN
				player.sendForAll(player.msg)
			}
			break
		}
	}

	fmt.Printf("Ok!\n\n")
	player.deck.setCardUsed(cardIt)

	player.msg.msgType = NEXT
	player.msg.cards[player.positionInIds] = player.deck.cards[cardIt]
	next := (player.positionInIds + 1) % NUM_PLAYERS
	player.ringClient.Send(player.clockWiseIds[next], player.msg)
}

/* Should be called by round master */
func (player *Player) WaitForAllCards() {
	player.ringClient.Recv(&player.msg)
	if player.msg.msgType == HEARTS_BROKEN {
		player.isHeartsBroken = true
		player.ringClient.Recv(&player.msg)
	}
	if player.msg.msgType != NEXT {
		player.msg.printUnexpectedType()
	}
}

/* Should be called by round master */
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
	return player.deck.cardsLeft == 0
}

//===================================================================
// Local functions
//===================================================================

func (p *Player) sendForAll(something any) {
	for it := range p.clockWiseIds {
		if it == int(p.positionInIds) {
			continue
		}
		p.ringClient.Send(p.clockWiseIds[it], something)
	}
}

func (d *deck) initDeck(myCards [MAX_CARDS_PER_ROUND]card) {
	d.cards = myCards
	d.cardsLeft = MAX_CARDS_PER_ROUND
	for i := range d.wasUsed {
		d.wasUsed[i] = false
	}
}

func (d *deck) getCardFromDeck(idx int) int {
	j := 1 /* first index is 1 */
	for i := range d.wasUsed {
		if !d.wasUsed[i] {
			if j == idx {
				return i
			}
			j++
		}
	}
	return -1
}

func (d *deck) setCardUsed(idx int) {
	d.wasUsed[idx] = true
	d.cardsLeft--
}

func (c *card) isSuitEqual(suit byte) bool {
	return (c.suit == suit)
}

func (c *card) isCardEqual(rank, suit byte) bool {
	return (c.rank == rank && c.suit == suit)
}

func (d *deck) printDeck() {
	j := 1 /* first index is 1 */
	fmt.Print("Your cards: ")
	for i := range d.wasUsed {
		if !d.wasUsed[i] {
			fmt.Printf("%v: %v%v ", j, ranks[d.cards[i].rank], suits[d.cards[i].suit])
			j++
		}
	}
	fmt.Printf("\n")
}

func (m *message) printUnexpectedType() {
	println("Message type not expected [%v]", m.msgType)
}
