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

func Start(players int) chan PlayerInfo{
	connectQueue := make(chan PlayerInfo)
	go start(players, connectQueue)
	return connectQueue
}

func start(maxPlayers int, queueConnect chan PlayerInfo){
	fmt.Printf("Waiting for players\n")
	players := make([]PlayerInfo, maxPlayers, maxPlayers)
	playersCount := 0
	for(playersCount < maxPlayers){
		player := <- queueConnect
		fmt.Printf("Player %s connected\n", player.Name)
		players[playersCount] = player
		playersCount++
	}

	fmt.Printf("Starting the game\n")
	ticker := time.NewTicker(time.Millisecond * 500)
	go func(){
		for {			
			<- ticker.C
			tick(players)
		}
	}()	
}

func tick(players []PlayerInfo){
	fmt.Printf("Tick\n")
	// check eaten 
	// spawn food
	// send stats to players
	for _, player := range players {
		fmt.Printf("Sending turn info to %s\n", player.Name)
		player.TurnSummaryCh <- TurnSummary { Round: 1}
	}
}
