package TokenRing

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

const MAX_CARDS_PER_ROUND = 13
const TOTAL_CARDS = 52
const NUM_PLAYERS = 4
const MAX_POINTS = 50
const HEARTS_VAL = 1
const QUEEN_SPADES_VAL = 13
const CARD_INVALID = -1

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
	Suit uint8
	Rank uint8
}

/* Message types */
const (
	CARDS = iota
	NEXT
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
	Cards          [MAX_CARDS_PER_ROUND]Card
	MsgType        uint8
	NumPlayedCards uint8
	EarnedPoints   uint8
	SourceId       uint8
}

type deck struct {
	cards     [MAX_CARDS_PER_ROUND]Card
	wasUsed   [MAX_CARDS_PER_ROUND]bool
	cardsLeft uint8
}

type Player struct {
	ringClient     TokenRingClient
	clockWiseIds   []uint8
	positionInIds  uint8
	myId           uint8
	deck           deck
	msg            Message
	points         uint8
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
		player.positionInIds = 0
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
		for i := byte(1); i < NUM_PLAYERS; i++ {
			if player.clockWiseIds[i] == myId {
				player.positionInIds = i
				break
			}
		}
	}
	player.myId = player.clockWiseIds[player.positionInIds]
	player.isCardDealer = isRingCreator
	player.points = 0
	player.deck.cardsLeft = 0
	player.isRoundMaster = false
	player.isGameActive = true
	fmt.Println("DEBUG: ids =", player.clockWiseIds)
	fmt.Println("DEBUG: pos =", player.positionInIds)
}

/* Should be called by card Dealer */
func (player *Player) DealCards() {
	var numbers []int
	for i := 0; i < TOTAL_CARDS; i++ {
		numbers = append(numbers, i)
	}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	fmt.Println("Cards shuffled! Sending cards!")

	var idNextRoundHead uint8
	pos := player.positionInIds
	for i := 0; i < TOTAL_CARDS; i++ {
		if (i%MAX_CARDS_PER_ROUND == 0) && (i != 0) {
			if pos == player.positionInIds { /* dealer cards */
				player.deck.initDeck(player.msg.Cards)
			} else { /* other players' cards */
				player.msg.MsgType = CARDS
				player.ringClient.Send(player.clockWiseIds[pos], &player.msg)
				fmt.Println("DEBUG: cards sent to", player.clockWiseIds[pos])
			}
			pos = (pos + 1) % NUM_PLAYERS
		}

		j := i % MAX_CARDS_PER_ROUND
		player.msg.Cards[j].Suit = uint8(numbers[i] / len(ranks))
		player.msg.Cards[j].Rank = uint8(numbers[i] % len(ranks))
		if player.msg.Cards[j].isCardEqual(TWO, CLUBS) {
			idNextRoundHead = player.clockWiseIds[pos]
		}
	}

	player.isHeartsBroken = false
	player.isRoundMaster = (idNextRoundHead == player.myId)
	if !player.isRoundMaster {
		player.msg.MsgType = BEGIN_ROUND
		player.ringClient.Send(idNextRoundHead, player.msg)
	}
}

/* Should be called by players (except Dealer) */
func (player *Player) GetCards() {
	player.isHeartsBroken = false
	fmt.Println("Getting cards...")
	player.ringClient.Recv(&player.msg)
	player.msg.inspectType(CARDS)
	fmt.Println("Got cards!")
	player.deck.initDeck(player.msg.Cards)
	for i := 0; i < MAX_CARDS_PER_ROUND; i++ {
		if player.deck.cards[i].isCardEqual(TWO, CLUBS) {
			player.isRoundMaster = true
			break
		}
	}
	if player.isRoundMaster {
		player.ringClient.Recv(&player.msg)
		player.msg.inspectType(BEGIN_ROUND)
	}
}

/* Show cards to player and let they choose */
func (player *Player) Play() {
	var selected int
	var cardIt int
	switch player.isRoundMaster {
	case true:
		/* reset number of played cards */
		player.msg.NumPlayedCards = 0
		fmt.Println("You begin round!")
		player.deck.printDeck()
		for {
			fmt.Print("Choose a card from your deck: ")
			fmt.Scanln(&selected)
			cardIt = player.deck.getCardFromDeck(selected)
			if cardIt == CARD_INVALID {
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
		if player.msg.MsgType == HEARTS_BROKEN {
			player.isHeartsBroken = true
			player.ringClient.Recv(&player.msg)
		}
		player.msg.inspectType(NEXT)
		masterSuit := player.getMasterSuit()
		hasMasterSuit := player.deckHasMasterSuit(masterSuit)
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
			/* broadcast HEARTS_BROKEN */
			if !player.isHeartsBroken && player.deck.cards[cardIt].isSuitEqual(HEARTS) {
				player.msg.MsgType = HEARTS_BROKEN
				player.sendForAll(&player.msg)
			}
			break
		}
	}

	fmt.Println("Ok!")
	player.deck.setCardUsed(cardIt)
	player.sendCardsToNextPlayer(cardIt)
}

/* Should be called by round master */
func (player *Player) WaitForAllCards() {
	player.ringClient.Recv(&player.msg)
	if player.msg.MsgType == HEARTS_BROKEN {
		player.isHeartsBroken = true
		player.ringClient.Recv(&player.msg)
	}
	player.msg.inspectType(NEXT)
}

/* Should be called by round master */
func (player *Player) InformRoundLoser() {
	/* reset isRoundMaster */
	player.isRoundMaster = false

	sum := uint8(0)
	/* obtain sum of points */
	for i := 0; i < NUM_PLAYERS; i++ {
		if player.msg.Cards[i].isSuitEqual(HEARTS) {
			sum += HEARTS_VAL
		} else if player.msg.Cards[i].isCardEqual(QUEEN, SPADES) {
			sum += QUEEN_SPADES_VAL
		}
	}
	if sum == 0 {
		fmt.Println("No one got points!")
		player.msg.MsgType = CONTINUE_GAME
		player.sendForAll(&player.msg)
		return
	}

	masterSuit := player.getMasterSuit()
	loserPosInIds := 0
	/* obtain loser id */
	for i := 0; i < NUM_PLAYERS; i++ {
		if player.msg.Cards[i].isSuitEqual(masterSuit) {
			if player.msg.Cards[i].Rank > player.msg.Cards[loserPosInIds].Rank {
				loserPosInIds = i
			}
		}
	}

	if player.clockWiseIds[loserPosInIds] != player.myId {
		player.msg.MsgType = ROUND_LOSER
		player.msg.EarnedPoints = sum
		player.msg.SourceId = player.myId
		player.ringClient.Send(player.clockWiseIds[loserPosInIds], &player.msg)
		/* wait for CONTINUE_GAME or MAX_PTS_REACHED */
		player.ringClient.Recv(&player.msg)
		if player.msg.MsgType == CONTINUE_GAME {
			player.msg.MsgType = CONTINUE_GAME
			player.sendForAll(&player.msg)
		} else if player.msg.MsgType == MAX_PTS_REACHED {
			player.isGameActive = false
		} else {
			println("Message type not expected [%v]", player.msg.MsgType)
		}
	} else {
		player.points += sum
		player.isRoundMaster = true
		fmt.Println("You lost round!")
		if player.points >= MAX_POINTS {
			player.isGameActive = false
			player.msg.MsgType = END_GAME
			player.sendForAll(&player.msg)
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
	minPts := player.points
	posInIds := player.positionInIds
	for i := uint8(0); i < NUM_PLAYERS; i++ {
		if i == player.positionInIds {
			continue
		}

		player.msg.MsgType = PTS_QUERY
		player.msg.SourceId = player.myId
		player.ringClient.Send(player.clockWiseIds[i], &player.msg)
		player.ringClient.Recv(&player.msg)
		player.msg.inspectType(PTS_REPLY)
		if player.msg.EarnedPoints < uint8(minPts) {
			posInIds = i
		}
	}
	player.msg.MsgType = GAME_WINNER
	player.ringClient.Send(player.clockWiseIds[posInIds], &player.msg)
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
		fmt.Println("Round master sent points query!")
		player.msg.MsgType = PTS_REPLY
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

func (p *Player) sendCardsToNextPlayer(cardIt int) {
	p.msg.MsgType = NEXT
	p.msg.NumPlayedCards++
	p.msg.Cards[p.positionInIds] = p.deck.cards[cardIt]
	next := (p.positionInIds + 1) % NUM_PLAYERS
	p.ringClient.Send(p.clockWiseIds[next], &p.msg)
}

func (p *Player) deckHasMasterSuit(masterSuit uint8) bool {
	for i := range p.deck.cards {
		if !p.deck.wasUsed[i] && p.deck.cards[i].isSuitEqual(masterSuit) {
			return true
		}
	}
	return false
}

func (p *Player) getMasterPosInIds() uint8 {
	return (p.positionInIds + (NUM_PLAYERS - p.msg.NumPlayedCards)) % NUM_PLAYERS
}

func (p *Player) getMasterSuit() uint8 {
	masterPos := p.getMasterPosInIds()
	log.Println("DEBUG: masterPos =", masterPos)
	masterSuit := p.msg.Cards[masterPos].Suit
	log.Println("DEBUG: masterSuit =", masterSuit)
	return masterSuit
}

func (p *Player) sendForAll(something any) {
	for it := range p.clockWiseIds {
		if it == int(p.positionInIds) {
			continue
		}
		p.ringClient.Send(p.clockWiseIds[it], something)
	}
}

func (d *deck) initDeck(myCards [MAX_CARDS_PER_ROUND]Card) {
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
	return CARD_INVALID
}

func (d *deck) setCardUsed(idx int) {
	d.wasUsed[idx] = true
	d.cardsLeft--
}

func (c *Card) isSuitEqual(suit byte) bool {
	return (c.Suit == suit)
}

func (c *Card) isCardEqual(rank, suit byte) bool {
	return (c.Rank == rank && c.Suit == suit)
}

func (d *deck) printDeck() {
	j := 1 /* first index is 1 */
	fmt.Print("Your cards: ")
	for i := range d.wasUsed {
		if !d.wasUsed[i] {
			fmt.Printf("%v: %v%v ", j, ranks[d.cards[i].Rank], suits[d.cards[i].Suit])
			j++
		}
	}
	fmt.Printf("\n")
}

func (m *Message) inspectType(typeT uint8) {
	if m.MsgType != typeT {
		println("Message type not expected [%v]", m.MsgType)
		os.Exit(1)
	}
}
