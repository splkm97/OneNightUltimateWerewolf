package WF

import (
	"math/rand"
	"time"
)

type CardDeck struct {
	Cards []string
}

func NewCardDeck() *CardDeck {
	cards := []string{
		"늑대인간",
		"하수인",
		"예언자",
		"강도",
		"말썽쟁이",
		"무두장이",
		"마을주민",
		"마을주민",
		"불면증환자",
	}
	cd := CardDeck{cards}
	cd.ShuffleCards()
	return &cd
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
