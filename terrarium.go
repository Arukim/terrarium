package main

import (
	"fmt"
	"github.com/arukim/terrarium/terra"
	"log"
	"math/rand"
	"os"
	"time"
)

func InitLog() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}

func main() {
	InitLog()

	var game = terra.NewGame(4, 25, 250*time.Millisecond)
	connectQueue := game.Start()

	players := [...]string{"Ivan", "Drake", "Sussana", "NyanCat"}

	for _, playerName := range players {
		player := terra.PlayerInfo{Name: playerName}
		player.TurnSummaryCh = make(chan terra.TurnSummary)
		player.PlayerTurnCh = make(chan terra.PlayerTurn)
		connectQueue <- player
		go func() {
			//player logic
			for {
				tInfo := <-player.TurnSummaryCh
				log.Printf("%s: got %d turn info\n", player.Name, tInfo.Turn)
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				player.PlayerTurnCh <- terra.PlayerTurn{}
			}
		}()
	}

	var input string
	fmt.Scanln(&input)
}
