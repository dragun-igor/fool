package game

import (
	"math/rand"
	"time"
)

type CardItem struct {
	ID           int    `json:"id"`
	Denomination string `json:"denomination"`
	Suit         string `json:"suit"`
	TrumpSuit    bool   `json:"trump_suit"`
	Selected     bool   `json:"selected"`
}

type DeckItem struct {
	Next *DeckItem
	Card CardItem
}

type Deck struct {
	FirstToken *DeckItem
	Length     int
}

type Pair struct {
	ID         int      `json:"id"`
	FirstCard  CardItem `json:"first_card"`
	SecondCard CardItem `json:"second_card"`
}

type Table struct {
	Pairs  []Pair `json:"pairs"`
	Putted int    `json:"-"`
	Cover  int    `json:"-"`
}

var suits = []string{"spades", "hearts", "diamonds", "clubs"}
var denominations = []string{"six", "seven", "eight", "nine", "ten", "jack", "queen", "king", "ace"}

func NewDeck() *Deck {
	cards := make([]CardItem, 0, 36)
	rand.Seed(time.Now().Unix())
	trumpSuit := suits[rand.Intn(4)]
	deck := Deck{}
	id := 0
	for _, suit := range suits {
		for _, denomination := range denominations {
			id++
			cards = append(cards, CardItem{
				ID:           id,
				Denomination: denomination,
				Suit:         suit,
				TrumpSuit:    trumpSuit == suit,
				Selected:     false,
			})
		}
	}
	for i := 0; i < 50; i++ {
		rand.Shuffle(36, func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	}
	for _, v := range cards {
		deck.AddCard(v)
	}
	return &deck
}

func (d *Deck) AddCard(card CardItem) {
	d.FirstToken = &DeckItem{
		Next: d.FirstToken,
		Card: card,
	}
	d.Length++
}

func (d *Deck) Get() CardItem {
	res := d.FirstToken.Card
	d.FirstToken = d.FirstToken.Next
	d.Length--
	return res
}

func (t *Table) PutCard(card CardItem) {
	card.Selected = false
	t.Pairs = append(t.Pairs, Pair{
		ID:        len(t.Pairs) + 1,
		FirstCard: card,
		SecondCard: CardItem{
			ID:           100,
			Denomination: "EXAMPLE",
			Suit:         "EXAMPLE",
		},
	})
	t.Putted++
}

func BringToHand(hand map[int]*CardItem, card CardItem) {
	hand[card.ID] = &card
}

func RemoveFromHand(hand map[int]*CardItem, id int) {
	delete(hand, id)
}
