package terra

type PlayerTurn struct {
}

type Player struct {
	Name          string
	TurnSummaryCh chan TurnSummary
	PlayerTurnCh  chan PlayerTurn
}

//player constructor
func NewPlayer(name string) *Player {
	p := new(Player)
	p.Name = name
	p.TurnSummaryCh = make(chan TurnSummary)
	p.PlayerTurnCh = make(chan PlayerTurn)
	return p
}
