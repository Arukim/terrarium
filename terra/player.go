package terra

import (
	"math/rand"
)

var _possibleNames = [...]string{"Ivan", "Drake", "Sussana", "NyanCat", "David Blain", "Sub-Zero", "Yakobovich"}

type PlayerTurn struct {
}

type Player struct {
	Name          string
	TurnSummaryCh chan TurnSummary
	PlayerTurnCh  chan PlayerTurn
}

//player constructor
func NewPlayer() *Player {
	p := new(Player)
	p.Name = _possibleNames[rand.Intn(len(_possibleNames))]
	p.TurnSummaryCh = make(chan TurnSummary)
	p.PlayerTurnCh = make(chan PlayerTurn)
	return p
}
