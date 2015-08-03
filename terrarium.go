package main

import (
	"fmt"
	"github.com/arukim/terrarium/terra"
)

func main() {
	connectQueue := terra.Start()

	players := [...]string{"Ivan","Drake","Sussana","NyanCat"}

	for _, playerName := range players {
		connectQueue <- terra.PlayerInfo{Name: playerName}
	}
	
	var input string
	fmt.Scanln(&input)
}
