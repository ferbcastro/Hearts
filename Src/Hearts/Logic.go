package Hearts

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

var cards = []byte{
	(two << 4) + (clubs),
	(three << 4) + (clubs),
	(four << 4) + (clubs),
	(five << 4) + (clubs),
	(six << 4) + (clubs),
	(seven << 4) + (clubs),
	(eight << 4) + (clubs),
	(nine << 4) + (clubs),
	(ten << 4) + (clubs),
	(jack << 4) + (clubs),
	(queen << 4) + (clubs),
	(king << 4) + (clubs),
	(ace << 4) + (clubs),
	(two << 4) + (diamonds),
	(three << 4) + (diamonds),
	(four << 4) + (diamonds),
	(five << 4) + (diamonds),
	(six << 4) + (diamonds),
	(seven << 4) + (diamonds),
	(eight << 4) + (diamonds),
	(nine << 4) + (diamonds),
	(ten << 4) + (diamonds),
	(jack << 4) + (diamonds),
	(queen << 4) + (diamonds),
	(king << 4) + (diamonds),
	(ace << 4) + (diamonds),
	(two << 4) + (spades),
	(three << 4) + (spades),
	(four << 4) + (spades),
	(five << 4) + (spades),
	(six << 4) + (spades),
	(seven << 4) + (spades),
	(eight << 4) + (spades),
	(nine << 4) + (spades),
	(ten << 4) + (spades),
	(jack << 4) + (spades),
	(queen << 4) + (spades),
	(king << 4) + (spades),
	(ace << 4) + (spades),
	(two << 4) + (hearts),
	(three << 4) + (hearts),
	(four << 4) + (hearts),
	(five << 4) + (hearts),
	(six << 4) + (hearts),
	(seven << 4) + (hearts),
	(eight << 4) + (hearts),
	(nine << 4) + (hearts),
	(ten << 4) + (hearts),
	(jack << 4) + (hearts),
	(queen << 4) + (hearts),
	(king << 4) + (hearts),
	(ace << 4) + (hearts),
}

type Player struct {
	myCards     []byte
	myMachineId byte
}

func printCard(rank int, suit int) {

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
