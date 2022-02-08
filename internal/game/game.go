package game

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

func (d *Deck) CardDistribution() CardMarks {
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
