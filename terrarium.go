package main

import (
	"fmt"
	"github.com/arukim/terrarium/terra"
)

func main() {
	var game = terra.NewGame(4,25)
	connectQueue := game.Start()

	players := [...]string{"Ivan","Drake","Sussana","NyanCat"}

	for _, playerName := range players {
		player := terra.PlayerInfo{Name: playerName}
		player.TurnSummaryCh = make(chan terra.TurnSummary)
		connectQueue <- player
		go func(){
			//player logic
			for{
				tInfo := <- player.TurnSummaryCh
				fmt.Printf("%s: got %d turn info\n",player.Name, tInfo.Turn)
			}
		}()
	}
	
	var input string
	fmt.Scanln(&input)
}
