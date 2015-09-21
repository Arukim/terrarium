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

	var game = terra.NewGame(4, 25, 25*time.Millisecond)
	connectQueue := game.Start()

	for i := 0; i < 4; i++ {
		player := terra.NewPlayer()
		connectQueue <- player
		go func() {
			//player logic
			for {
				tInfo := <-player.TurnSummaryCh
				log.Printf("%s: got %d turn info\n", player.Name, tInfo.Turn)
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				player.PlayerTurnCh <- terra.PlayerTurn{}
			}
		}()
	}

	var input string
	fmt.Scanln(&input)
}
