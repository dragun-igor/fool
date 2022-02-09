package game

import (
	"math/rand"
	"time"
)

type CardMarks struct {
	Denomination string `json:"denomination"`
	Suit         string `json:"suit"`
}

type Hand []CardMarks

type PutCard struct {
	Card CardMarks `json:"card"`
}

type DeckToken struct {
	Next *DeckToken
	Card CardMarks
}

type Deck struct {
	FirstToken  *DeckToken
	Length      int
	PrimarySuit string
}

var suits = []string{"spades", "hearts", "diamonds", "clubs"}
var denominations = []string{"six", "seven", "eight", "nine", "ten", "jack", "queen", "king", "ace"}

func NewDeck() *Deck {
	rand.Seed(time.Now().Unix())
	deck := Deck{
		PrimarySuit: suits[rand.Intn(4)],
	}
	for _, suit := range suits {
		for _, denomination := range denominations {
			deck.AddCard(denomination, suit)
		}
	}
	return &deck
}

func (d *Deck) AddCard(denomination string, suit string) {
	d.FirstToken = &DeckToken{
		Next: d.FirstToken,
		Card: CardMarks{
			Denomination: denomination,
			Suit:         suit,
		},
	}
	d.Length++
}

func (d *Deck) Get() CardMarks {
	res := d.FirstToken.Card
	d.FirstToken = d.FirstToken.Next
	d.Length--
	return res
}

func BringToHand(hand Hand, cards ...CardMarks) Hand {
	return append(hand, cards...)
}

func RemoveFromHand(hand Hand, card CardMarks) Hand {
	for i := range hand {
		if hand[i] == card {
			hand = append(hand[:i], hand[i+1:]...)
			break
		}
	}
	return hand
}
