package game

import (
	"math/rand"
	"time"
)

type CardItem struct {
	Denomination string `json:"denomination"`
	Suit         string `json:"suit"`
}

type Hand []CardItem

type DeckItem struct {
	Next *DeckItem
	Card CardItem
}

type Deck struct {
	FirstToken *DeckItem
	Length     int
	TrumpSuit  string
}

var suits = []string{"spades", "hearts", "diamonds", "clubs"}
var denominations = []string{"six", "seven", "eight", "nine", "ten", "jack", "queen", "king", "ace"}

func NewDeck() *Deck {
	rand.Seed(time.Now().Unix())
	deck := Deck{
		TrumpSuit: suits[rand.Intn(4)],
	}
	for _, suit := range suits {
		for _, denomination := range denominations {
			deck.AddCard(denomination, suit)
		}
	}
	return &deck
}

func (d *Deck) AddCard(denomination string, suit string) {
	d.FirstToken = &DeckItem{
		Next: d.FirstToken,
		Card: CardItem{
			Denomination: denomination,
			Suit:         suit,
		},
	}
	d.Length++
}

func (d *Deck) Get() CardItem {
	res := d.FirstToken.Card
	d.FirstToken = d.FirstToken.Next
	d.Length--
	return res
}

func BringToHand(hand Hand, cards ...CardItem) Hand {
	return append(hand, cards...)
}

func RemoveFromHand(hand Hand, card CardItem) Hand {
	for i := range hand {
		if hand[i] == card {
			hand = append(hand[:i], hand[i+1:]...)
			break
		}
	}
	return hand
}
