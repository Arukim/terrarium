package main

import (
	"fmt"
	"github.com/arukim/terrarium/terra"
)

func main() {
	connectQueue := terra.Start(4)

	players := [...]string{"Ivan","Drake","Sussana","NyanCat"}

	for _, playerName := range players {
		player := terra.PlayerInfo{Name: playerName}
		player.TurnSummaryCh = make(chan terra.TurnSummary)
		connectQueue <- player
		go func(){
			//player logic
			for{
				<- player.TurnSummaryCh
			}
		}()
	}
	
	var input string
	fmt.Scanln(&input)
}
