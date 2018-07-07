package types

import (
	"testing"
	"github.com/stretchr/testify/assert"
)



func Test_parseCardsData(t *testing.T){
	cardsData := "(69C:46% )(69C:60% )(69C:50% )(69C:51% )(69C:55% )(69C:49% )(70C:64% )(69C:47% )(69C:29% )(69C:59% )(69C:48% )(69C:48% )"
	cards := parseCardsData(cardsData)
	assert.NotEmpty(t, cards)
	assert.Equal(t,69, cards[0].temp)
	assert.Equal(t,46, cards[0].load)
}
