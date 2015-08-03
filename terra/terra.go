package terra

import (
	"fmt"
	"time"
)

type PlayerInfo struct{
	Name string
}

func Start() chan PlayerInfo{
	connectQueue := make(chan PlayerInfo)
	go start(connectQueue)
	return connectQueue
}

func start(queueConnect chan PlayerInfo){
	fmt.Printf("Waiting for players\n")
	players := make([]PlayerInfo, 4, 4)
	playersCount := 0
	for(playersCount < 4){
		player := <- queueConnect
		fmt.Printf("Player %s connected\n", player)
		players[playersCount] = player
		playersCount++
	}
	
	ticker := time.NewTicker(time.Millisecond * 500)
	go func(){
		for {
			<- ticker.C
			go tick()
		}
	}()	
}

func tick(){
	fmt.Printf("Tick\n")
}
