package terra

import (
	"fmt"
	"time"
)

type TurnSummary struct {
	Turn int
}

type PlayerTurn struct {
}

type PlayerInfo struct{
	Name string
	TurnSummaryCh chan TurnSummary
	PlayerTurnCh chan PlayerTurn	
}

type Game struct {
	// settings
	maxPlayers int
	maxTurns int
	// current info
	players []PlayerInfo
	playersCount int
	turn int
}

func NewGame(maxPlayers int, maxTurns int) *Game {
	g := new(Game)
	g.maxPlayers = maxPlayers
	g.maxTurns = maxTurns
	return g
}

func (g *Game) Start() chan PlayerInfo{
	connectQueue := make(chan PlayerInfo)
	go func(){
		fmt.Printf("Waiting for players\n")
		g.players = make([]PlayerInfo, g.maxPlayers, g.maxPlayers)
		g.playersCount = 0
		for(g.playersCount < g.maxPlayers){
			player := <- connectQueue
			fmt.Printf("Player %s connected\n", player.Name)
			g.players[g.playersCount] = player
			g.playersCount++
		}
		
		fmt.Printf("Starting the game\n")
		ticker := time.NewTicker(time.Millisecond * 500)
		go func(){
			for g.turn < g.maxTurns {			
				<- ticker.C
				g.Tick()				
			}
			ticker.Stop()
			fmt.Printf("Game has ended\n")
			// check for winner
		}()	
	}()
	return connectQueue
}

func (g *Game) Tick(){
	fmt.Printf("Turn %d\n", g.turn)
	// check eaten 
	// spawn food
	// send stats to players
	for _, player := range g.players {
		go func(player PlayerInfo, turn int){
			player.TurnSummaryCh <- TurnSummary { Turn: turn}
		}(player, g.turn)
	}
	g.turn++
}
