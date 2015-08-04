package terra

import (
	"fmt"
	"time"
)

type TurnSummary struct {
	Round int
}

type PlayerTurn struct {
}

type PlayerInfo struct{
	Name string
	TurnSummaryCh chan TurnSummary
	PlayerTurnCh chan PlayerTurn	
}

type Game struct {
	Players []PlayerInfo
	PlayersCount int
	MaxPlayers int
}

func Start(players int) chan PlayerInfo{
	connectQueue := make(chan PlayerInfo)
	go start(players, connectQueue)
	return connectQueue
}

func start(maxPlayers int, queueConnect chan PlayerInfo){
	fmt.Printf("Waiting for players\n")
	game := Game{MaxPlayers: maxPlayers}
	game.Players = make([]PlayerInfo, game.MaxPlayers, game.MaxPlayers)
	game.PlayersCount = 0
	for(game.PlayersCount < game.MaxPlayers){
		player := <- queueConnect
		fmt.Printf("Player %s connected\n", player.Name)
		game.Players[game.PlayersCount] = player
		game.PlayersCount++
	}

	fmt.Printf("Starting the game\n")
	ticker := time.NewTicker(time.Millisecond * 500)
	go func(){
		for {			
			<- ticker.C
			game.Tick()
		}
	}()	
}

func (g Game) Tick(){
	fmt.Printf("Tick\n")
	// check eaten 
	// spawn food
	// send stats to players
	for _, player := range g.Players {
		fmt.Printf("Sending turn info to %s\n", player.Name)
		player.TurnSummaryCh <- TurnSummary { Round: 1}
	}
}
