package terra

import (
	"github.com/arukim/terrarium/helpers"
	"log"
	"time"
)

type TurnSummary struct {
	Turn int
}

type PlayerTurn struct {
}

type PlayerInfo struct {
	Name          string
	TurnSummaryCh chan TurnSummary
	PlayerTurnCh  chan PlayerTurn
}

type Game struct {
	// settings
	maxPlayers  int
	maxTurns    int
	turnTimeout time.Duration
	// current info
	players      []PlayerInfo
	playersCount int
	turn         int
}

func NewGame(maxPlayers int, maxTurns int, turnTimeout time.Duration) *Game {
	g := new(Game)
	g.maxPlayers = maxPlayers
	g.maxTurns = maxTurns
	g.turnTimeout = turnTimeout
	return g
}

func (g *Game) Start() chan PlayerInfo {
	connectQueue := make(chan PlayerInfo)
	go func() {
		log.Printf("Waiting for players\n")
		g.players = make([]PlayerInfo, g.maxPlayers, g.maxPlayers)
		g.playersCount = 0
		for g.playersCount < g.maxPlayers {
			player := <-connectQueue
			log.Printf("Player %s connected\n", player.Name)
			g.players[g.playersCount] = player
			g.playersCount++
		}

		log.Printf("Starting the game\n")
		go func() {
			for g.turn < g.maxTurns {
				g.Turn()
			}
			log.Printf("Game has ended\n")
			// check for winner
		}()
	}()
	return connectQueue
}

func (g *Game) Turn() {
	log.Printf("Turn %d\n", g.turn)
	// check eaten
	// spawn food
	// send stats to players
	for _, player := range g.players {
		// check for old turns
		select {
		case <-player.PlayerTurnCh:
		default:
		}

		// start timeout, send turn summary
		var timeout = helpers.NewTimeout(g.turnTimeout)
		player.TurnSummaryCh <- TurnSummary{Turn: g.turn}

		// wait for player answer
		select {
		case <-player.PlayerTurnCh:
			log.Printf("turn was made")
		case <-timeout.Alarm:
			log.Printf("turn timeout")
		}
	}
	g.turn++
}
