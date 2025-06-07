package Hearts

import (
	"Src/TokenRing"
	"fmt"
	"log"
	"math/rand"
)

const CARDS_PER_ROUND = 13
const TOTAL_CARDS = 52
const NUM_PLAYERS = 4
const MAX_POINTS = 50
const HEARTS_VAL = 1
const QUEEN_SPADES_VAL = 13

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

type Card struct {
	Suit int8
	Rank int8
}

/* Message types */
const (
	NEXT = iota
	PTS_QUERY
	PTS_REPLY
	GAME_WINNER
	END_GAME
	MAX_PTS_REACHED
	ROUND_LOSER
	CONTINUE_GAME
	HEARTS_BROKEN
	BEGIN_ROUND
)

type Message struct {
	Cards          [NUM_PLAYERS]Card
	MsgType        byte
	NumPlayedCards byte
	SourceId       byte
	MasterSuit     int8
	EarnedPoints   int
}

type deck struct {
	cards     [CARDS_PER_ROUND]Card
	cardsLeft int
}

type Player struct {
	ringClient     TokenRing.TokenRingClient
	clockWiseIds   []byte
	myId           byte
	myPosition     int
	points         int
	deck           deck
	msg            Message
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
		fmt.Scanln(&ips[0])
		fmt.Println("Now enter other players' ip in clockwise order")
		for i := 1; i < NUM_PLAYERS; i++ {
			fmt.Print("Enter ip: ")
			fmt.Scanln(&ips[i])
		}
		/* creat ring */
		ids = player.ringClient.CreateRing(ips)
		player.clockWiseIds = ids
		player.myPosition = 0
		/* broadcast ids */
		player.sendForAll(&ids)
	} else {
		/* enter ring */
		fmt.Print("Enter your ip: ")
		fmt.Scanln(&ip)
		myId = byte(player.ringClient.EnterRing(ip))
		/* read ids */
		player.ringClient.Recv(&player.clockWiseIds)
		/* discover position */
		for i := 1; i < NUM_PLAYERS; i++ {
			if player.clockWiseIds[i] == myId {
				player.myPosition = i
				break
			}
		}
	}
	player.points = 0
	player.deck.cardsLeft = 0
	player.isRoundMaster = false
	player.isGameActive = true
	player.isCardDealer = isRingCreator
	player.myId = player.clockWiseIds[player.myPosition]
}

/* Should be called by card Dealer */
func (player *Player) DealCards() {
	var numbers []int8
	var idRoundMaster byte
	var cards [CARDS_PER_ROUND]Card

	for i := int8(0); i < TOTAL_CARDS; i++ {
		numbers = append(numbers, i)
	}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	fmt.Println("Cards shuffled! Sending cards!")

	numbersIt := 0
	for i := 0; i < NUM_PLAYERS; i++ {
		for j := 0; j < CARDS_PER_ROUND; j++ {
			cards[j].Rank = numbers[numbersIt] % int8(len(ranks))
			cards[j].Suit = numbers[numbersIt] / int8(len(ranks))
			if cards[j].isCardEqual(TWO, CLUBS) {
				idRoundMaster = player.clockWiseIds[i]
			}
			numbersIt++
		}
		if player.clockWiseIds[i] == player.myId {
			player.deck.initDeck(cards)
		} else {
			player.ringClient.Send(player.clockWiseIds[i], &cards)
		}
	}

	request := BEGIN_ROUND
	if idRoundMaster == player.myId {
		player.isRoundMaster = true
	} else {
		player.ringClient.Send(idRoundMaster, &request)
	}
}

/* Should be called by players (except Dealer) */
func (player *Player) GetCards() {
	var cards [CARDS_PER_ROUND]Card
	var request int

	fmt.Println("Getting cards...")
	player.ringClient.Recv(&cards)
	player.deck.initDeck(cards)
	fmt.Println("Got cards!")
	for i := range cards {
		if cards[i].isCardEqual(TWO, CLUBS) {
			player.isRoundMaster = true
			break
		}
	}

	if !player.isRoundMaster {
		return
	}

	for {
		player.ringClient.Recv(&request)
		if request == BEGIN_ROUND {
			break
		}
	}
}

/* Show cards to player and let they choose */
func (player *Player) Play() {
	var selected int
	var card *Card
	var hasMasterSuit bool
	var masterSuit int8
	var next int

	switch player.isRoundMaster {
	case true:
		player.msg.NumPlayedCards = 0
		fmt.Println("You begin round!")
		player.deck.printDeck()
		for {
			fmt.Print("Choose a card from your deck: ")
			fmt.Scanln(&selected)
			if card = player.deck.getCardFromDeck(selected - 1); card == nil {
				fmt.Println("Invalid card!")
				continue
			}
			log.Println("DEBUG:", card)
			if card.isSuitEqual(HEARTS) && !player.isHeartsBroken {
				fmt.Println("Invalid card! Hearts not broken!")
				continue
			}
			break
		}
		player.msg.MasterSuit = card.Suit
	case false:
		for {
			player.ringClient.Recv(&player.msg)
			if player.msg.MsgType == HEARTS_BROKEN {
				player.SetHeartsBroken()
			}
			if player.msg.MsgType == NEXT {
				break
			}
		}
		fmt.Println("Your turn!")
		player.deck.printDeck()
		player.printRecvCards()

		masterSuit = player.msg.MasterSuit
		hasMasterSuit = player.deckHasMasterSuit(masterSuit)
		for {
			fmt.Print("Choose a card from your deck: ")
			fmt.Scanln(&selected)
			if card = player.deck.getCardFromDeck(selected - 1); card == nil {
				fmt.Println("Invalid card!")
				continue
			}
			if hasMasterSuit && !card.isSuitEqual(masterSuit) {
				fmt.Println("Invalid card! Suit not the same as master suit!")
				continue
			}
			/* broadcast HEARTS_BROKEN */
			if !player.isHeartsBroken && card.isSuitEqual(HEARTS) {
				player.SetHeartsBroken()
				player.msg.MsgType = HEARTS_BROKEN
				player.sendForAll(&player.msg)
			}
			break
		}
	}

	fmt.Printf("Ok!\n\n")
	player.deck.setCardUsed(selected)

	player.msg.MsgType = NEXT
	player.msg.NumPlayedCards++
	player.msg.Cards[player.myPosition] = *card
	next = (player.myPosition + 1) % NUM_PLAYERS
	player.ringClient.Send(player.clockWiseIds[next], &player.msg)
}

func (player *Player) SetHeartsBroken() {
	player.isHeartsBroken = true
	fmt.Println("Hearts broken!")
}

func (player *Player) ResetHeartsBroken() {
	player.isHeartsBroken = false
}

/* Should be called by round master */
func (player *Player) WaitForAllCards() {
	for {
		player.ringClient.Recv(&player.msg)
		if player.msg.MsgType == HEARTS_BROKEN {
			player.SetHeartsBroken()
		}
		/* player before round master makes no distinction between master and other
		 * players and that is why message type here should be NEXT */
		if player.msg.MsgType == NEXT {
			break
		}
	}
}

/* Should be called by round master */
func (player *Player) InformRoundLoser() {
	var loserCard Card
	var loserId byte

	/* reset isRoundMaster */
	player.isRoundMaster = false

	sum := 0
	/* obtain sum of points */
	for i := 0; i < NUM_PLAYERS; i++ {
		if player.msg.Cards[i].isSuitEqual(HEARTS) {
			sum += HEARTS_VAL
		} else if player.msg.Cards[i].isCardEqual(QUEEN, SPADES) {
			sum += QUEEN_SPADES_VAL
		}
	}
	if sum == 0 {
		fmt.Printf("No one got points!\n\n")
		player.msg.MsgType = CONTINUE_GAME
		player.sendForAll(&player.msg)
		return
	}

	masterSuit := player.msg.MasterSuit
	loserCard.Rank = -1
	/* obtain loser id */
	for i := 0; i < NUM_PLAYERS; i++ {
		if player.msg.Cards[i].isSuitEqual(masterSuit) {
			if player.msg.Cards[i].Rank > loserCard.Rank {
				loserCard = player.msg.Cards[i]
				loserId = player.clockWiseIds[i]
			}
		}
	}

	if loserId != player.myId {
		player.msg.MsgType = ROUND_LOSER
		player.msg.EarnedPoints = sum
		player.msg.SourceId = player.myId
		player.ringClient.Send(loserId, &player.msg)
		/* wait for CONTINUE_GAME or MAX_PTS_REACHED */
		for {
			player.ringClient.Recv(&player.msg)
			if player.msg.MsgType == CONTINUE_GAME {
				player.msg.MsgType = CONTINUE_GAME
				player.sendForAll(&player.msg)
				break
			} else if player.msg.MsgType == MAX_PTS_REACHED {
				player.isGameActive = false
				break
			}
		}
	} else {
		player.points += sum
		player.isRoundMaster = true
		fmt.Println("You lost round!")
		if player.points >= MAX_POINTS {
			player.isGameActive = false
		} else {
			player.msg.MsgType = CONTINUE_GAME
			player.sendForAll(&player.msg)
		}
	}
}

/* Should be called by round master */
func (player *Player) IsThereAWinner() bool {
	return !player.isGameActive
}

/* Round master sends a message to the round winner */
func (player *Player) AnounceWinner() {
	var idWinner byte
	var currentMin int
	var dest byte

	sentReplies := 0
	currentMin = player.points
	idWinner = player.myId
	for it := byte(0); it < NUM_PLAYERS; it++ {
		dest = player.clockWiseIds[it]
		if dest == player.myId {
			continue
		}

		player.msg.MsgType = PTS_QUERY
		player.msg.SourceId = player.myId
		player.ringClient.Send(dest, &player.msg)
		for {
			player.ringClient.Recv(&player.msg)
			if player.msg.MsgType == PTS_REPLY && player.msg.SourceId == dest {
				break
			}
		}
		if currentMin > player.msg.EarnedPoints {
			currentMin = player.msg.EarnedPoints
			idWinner = player.msg.SourceId
			sentReplies++
		}

		if sentReplies == NUM_PLAYERS {
			break
		}
	}

	fmt.Println("We have a winner!")

	player.msg.MsgType = GAME_WINNER
	player.ringClient.Send(idWinner, &player.msg)
	player.msg.MsgType = END_GAME
	player.sendForAll(&player.msg)
}

/* Players should call this (except round master) */
func (player *Player) WaitForResult() {
	player.ringClient.Recv(&player.msg)
	switch player.msg.MsgType {
	case CONTINUE_GAME:
		fmt.Println("New round!")
	case ROUND_LOSER:
		fmt.Println("You lost round!")
		player.isRoundMaster = true
		player.points += player.msg.EarnedPoints
		if player.points >= MAX_POINTS {
			player.msg.MsgType = MAX_PTS_REACHED
			player.ringClient.Send(player.msg.SourceId, &player.msg)
			player.WaitForResult() /* recursive call */
		} else {
			player.msg.MsgType = CONTINUE_GAME
			player.ringClient.Send(player.msg.SourceId, &player.msg)
			player.WaitForResult() /* recursive call */
		}
	case PTS_QUERY:
		player.msg.MsgType = PTS_REPLY
		player.msg.SourceId = player.myId
		player.msg.EarnedPoints = player.points
		player.ringClient.Send(player.msg.SourceId, &player.msg)
		player.WaitForResult() /* recursive call */
	case GAME_WINNER:
		fmt.Println("You won!")
		player.WaitForResult() /* recursive call */
	case END_GAME:
		player.isGameActive = false
	}
}

func (player *Player) PrintPoints() {
	fmt.Println("Your points:", player.points)
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

func (player *Player) deckHasMasterSuit(masterSuit int8) bool {
	for i := 0; i < player.deck.cardsLeft; i++ {
		if player.deck.cards[i].isSuitEqual(masterSuit) {
			return true
		}
	}
	return false
}

func (p *Player) sendForAll(something any) {
	for it := range p.clockWiseIds {
		if it == p.myPosition {
			continue
		}
		p.ringClient.Send(p.clockWiseIds[it], something)
	}
}

func (d *deck) initDeck(myCards [CARDS_PER_ROUND]Card) {
	d.cards = myCards
	d.cardsLeft = CARDS_PER_ROUND
}

func (deck *deck) getCardFromDeck(idx int) *Card {
	if idx >= deck.cardsLeft {
		return nil
	}
	return &deck.cards[idx]
}

func (deck *deck) setCardUsed(idx int) {
	deck.cards[idx] = deck.cards[deck.cardsLeft-1]
	deck.cardsLeft--
}

func (c *Card) isSuitEqual(suit int8) bool {
	return (c.Suit == suit)
}

func (c *Card) isCardEqual(rank, suit int8) bool {
	return (c.Rank == rank && c.Suit == suit)
}

func (card *Card) printCard(it int) {
	//fmt.Println("DEBUG:", card.Rank, card.Suit)
	fmt.Printf("%v: %v%v ", it, ranks[card.Rank], suits[card.Suit])
}

func (player *Player) printRecvCards() {
	numPlayed := int(player.msg.NumPlayedCards)
	it := (player.myPosition + NUM_PLAYERS - numPlayed) % NUM_PLAYERS
	log.Println("it", it, "numPlayed", numPlayed)
	fmt.Println("Cards played:")
	for {
		player.msg.Cards[it].printCard(int(it))
		if it = (it + 1) % NUM_PLAYERS; it == player.myPosition {
			break
		}
	}
	fmt.Println()
}

func (deck *deck) printDeck() {
	fmt.Println("Your cards:")
	for i := range deck.cards {
		deck.cards[i].printCard(i + 1)
		if i%7 == 0 && i != 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
