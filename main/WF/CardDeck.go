package WF

import (
	"math/rand"
	"sort"
	"time"
)

type CardDeck struct {
	Cards      []string
	ChoiceChan chan int
}

func NewCardDeck() *CardDeck {
	cd := CardDeck{
		Cards:      make([]string, 0, 10),
		ChoiceChan: make(chan int, 2),
	}
	cd.ShuffleCards()
	return &cd
}
func (cd *CardDeck) SortCards() {
	sort.Strings(cd.Cards)
}
func (cd *CardDeck) ShuffleCards() {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(cd.Cards), func(i, j int) {
		cd.Cards[i], cd.Cards[j] = cd.Cards[j], cd.Cards[i]
	})
}

func (cd *CardDeck) PopCard() string {
	temp := cd.Cards[0]
	cd.Cards = cd.Cards[1:]
	return temp
}
