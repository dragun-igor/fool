package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Deck struct {
	FirstToken *DeckItem
	Len        int
}

type DeckItem struct {
	Next *DeckItem
	Card CardItem
}

type CardItem struct {
	ID           int    `json:"id"`
	Denomination string `json:"denomination"`
	Suit         string `json:"suit"`
	TrumpSuit    bool   `json:"trump_suit"`
	Selected     bool   `json:"selected"`
	Help         bool   `json:"help"`
}

type Table struct {
	Pairs   []Pair `json:"pairs"`
	Putted  int    `json:"-"`
	Covered int    `json:"-"`
}

type Pair struct {
	ID         int      `json:"id"`
	FirstCard  CardItem `json:"first_card"`
	SecondCard CardItem `json:"second_card"`
	Covered    bool     `json:"-"`
}

type Hand struct {
	Defend         bool
	Hand           map[int]CardItem
	SelectedCardID int
	HelpID         []int
}

var (
	suits              = [4]string{"spades", "hearts", "diamonds", "clubs"}
	denominations      = [9]string{"six", "seven", "eight", "nine", "ten", "jack", "queen", "king", "ace"}
	denominationWeight = map[string]int{
		"six":   1,
		"seven": 2,
		"eight": 3,
		"nine":  4,
		"ten":   5,
		"jack":  6,
		"queen": 7,
		"king":  8,
		"ace":   9,
	}
)

func NewDeck() *Deck {
	cards := make([]CardItem, 0, 36)
	rand.Seed(time.Now().Unix())
	trumpSuit := suits[rand.Intn(4)]
	deck := &Deck{}
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
		rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	}
	for i := range cards {
		deck.Add(cards[i])
	}
	return deck
}

func NewTable() *Table {
	return &Table{
		Pairs:   make([]Pair, 0, 6),
		Putted:  0,
		Covered: 0,
	}
}

func StartNewGame() (*Deck, *Table) {
	return NewDeck(), NewTable()
}

func (d *Deck) Add(card CardItem) {
	d.FirstToken = &DeckItem{
		Next: d.FirstToken,
		Card: card,
	}
	d.Len++
}

func (d *Deck) Get(hand map[int]CardItem) error {
	var err error
	id := d.FirstToken.Card.ID
	if _, ok := hand[id]; !ok {
		hand[d.FirstToken.Card.ID] = d.FirstToken.Card
		d.FirstToken = d.FirstToken.Next
		d.Len--
	} else {
		err = fmt.Errorf("hand and deck cards have same id")
	}
	return err
}

func (d *Deck) GetHand() (Hand, error) {
	var err error
	hand := make(map[int]CardItem, 18)
	for i := 0; i < 6; i++ {
		err = d.Get(hand)
	}
	return Hand{Hand: hand}, err
}

func (t *Table) PutCardOnTable(hand Hand) (Hand, error) {
	if len(t.Pairs) >= 6 {
		return hand, fmt.Errorf("too much cards on table")
	}
	if card, ok := hand.Hand[hand.SelectedCardID]; !ok {
		return hand, fmt.Errorf("selected card isn't in hand")
	} else {
		if ok := t.CheckPutCard(card); !ok {
			// do nothing
		} else {
			card.Selected = false
			id := len(t.Pairs) + 1
			t.Pairs = append(t.Pairs, Pair{
				ID:        id,
				FirstCard: card,
			})
			t.Putted++
			delete(hand.Hand, hand.SelectedCardID)
			hand.SelectedCardID = 0
		}
	}
	return hand, nil
}

func (t *Table) CoverCardOnTable(id int, hand Hand) (Hand, error) {
	card, ok := hand.Hand[hand.SelectedCardID]
	if !ok {
		return hand, fmt.Errorf("error")
	}
	if ok := t.CheckCoverCard(id, card); !ok {
		return hand, fmt.Errorf("ne kroet")
	}
	delete(hand.Hand, hand.SelectedCardID)
	card.Selected = false
	hand.SelectedCardID = 0
	t.Pairs[id-1].SecondCard = card
	t.Pairs[id-1].Covered = true
	t.Covered++
	return hand, nil
}

func (t *Table) BringCardsToHand(hand map[int]CardItem) {
	for i := range t.Pairs {
		hand[t.Pairs[i].FirstCard.ID] = t.Pairs[i].FirstCard
		if !t.Pairs[i].Covered {
			// do nothing
		} else {
			hand[t.Pairs[i].SecondCard.ID] = t.Pairs[i].SecondCard
		}
	}
	t.Clear()
}

func (t *Table) Clear() {
	t.Pairs = make([]Pair, 0, 6)
	t.Putted = 0
	t.Covered = 0
}

func selectOff(hand Hand) (Hand, error) {
	if card, ok := hand.Hand[hand.SelectedCardID]; !ok {
		return hand, fmt.Errorf("selected card not in hand")
	} else {
		card.Selected = false
		hand.Hand[hand.SelectedCardID] = card
		hand.SelectedCardID = 0
		return hand, nil
	}
}

func selectOn(id int, hand Hand) (Hand, error) {
	if card, ok := hand.Hand[id]; !ok {
		return hand, fmt.Errorf("selected card not in hand")
	} else {
		card.Selected = true
		hand.Hand[id] = card
		hand.SelectedCardID = id
		return hand, nil
	}
}

//SelectCard changes card state
func (t *Table) SelectCard(id int, hand Hand) (Hand, error) {
	var err error
	if hand.SelectedCardID != 0 {
		if hand.SelectedCardID != id {
			hand, err = selectOff(hand)
			hand, err = selectOn(id, hand)
		} else {
			hand, err = selectOff(hand)
		}
	} else {
		hand, err = selectOn(id, hand)
	}
	if err != nil {
		log.Println(err)
	}
	if !hand.Defend {
		// do nothing
	} else {
		t.SelectedCardCanCover(hand)
	}
	return hand, nil
}

// SelectedCardCanCover changes state of cards for indicate card can cover other card on table
func (t *Table) SelectedCardCanCover(hand Hand) {
	card := hand.Hand[hand.SelectedCardID]
	for i := range t.Pairs {
		firstCardPair := t.Pairs[i].FirstCard

		if !t.Pairs[i].Covered {
			// do nothing
		} else {
			t.Pairs[i].FirstCard.Selected = false
			continue
		}

		if card.TrumpSuit && !firstCardPair.TrumpSuit {
			t.Pairs[i].FirstCard.Selected = true
			continue
		}

		if card.TrumpSuit &&
			firstCardPair.TrumpSuit &&
			denominationWeight[card.Denomination] > denominationWeight[firstCardPair.Denomination] {

			t.Pairs[i].FirstCard.Selected = true
			continue
		}

		if card.Suit == firstCardPair.Suit &&
			denominationWeight[card.Denomination] > denominationWeight[firstCardPair.Denomination] {

			t.Pairs[i].FirstCard.Selected = true
			continue
		}

		t.Pairs[i].FirstCard.Selected = false
	}
}

func (t Table) HelperCardCanPut(hand Hand) (Hand, error) {
	den := make([]string, 0, 6)
	for i := range t.Pairs {
		den = append(den, t.Pairs[i].FirstCard.Denomination)
		if !t.Pairs[i].Covered {
			// do nothing
		} else {
			den = append(den, t.Pairs[i].SecondCard.Denomination)
		}
	}
	var temp bool
	for key := range hand.Hand {
		temp = false
		for i := range den {
			if den[i] == hand.Hand[key].Denomination {
				temp = true
			}
		}
		if card, ok := hand.Hand[key]; !ok {
			// do nothing
		} else {
			card.Help = temp
			hand.Hand[key] = card
		}
	}
	return hand, nil
}

//SelectedCardCanCoverClear clears helper
func (t *Table) SelectedCardCanCoverClear() {
	for i := range t.Pairs {
		t.Pairs[i].FirstCard.Selected = false
		t.Pairs[i].FirstCard.Help = false
		if !t.Pairs[i].Covered {
			// do nothing
		} else {
			t.Pairs[i].SecondCard.Selected = false
			t.Pairs[i].SecondCard.Help = false
		}
	}
}

func (t Table) CheckPutCard(card CardItem) bool {
	if len(t.Pairs) == 0 {
		return true
	}
	denomination := card.Denomination
	for i := range t.Pairs {
		if t.Pairs[i].FirstCard.Denomination == denomination ||
			t.Pairs[i].SecondCard.Denomination == denomination {

			return true
		}
	}
	return false
}

func (t Table) CheckCoverCard(id int, secondCard CardItem) bool {
	pair := t.Pairs[id-1]
	if !pair.Covered {
		// do nothing
	} else {
		return false
	}
	if !pair.FirstCard.TrumpSuit && secondCard.TrumpSuit {
		return true
	}
	if pair.FirstCard.TrumpSuit &&
		secondCard.TrumpSuit &&
		denominationWeight[pair.FirstCard.Denomination] < denominationWeight[secondCard.Denomination] {

		return true
	}

	if pair.FirstCard.Suit == secondCard.Suit &&
		denominationWeight[pair.FirstCard.Denomination] < denominationWeight[secondCard.Denomination] {

		return true
	}

	return false
}
