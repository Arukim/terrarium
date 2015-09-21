package main

import (
	"fmt"
	"github.com/arukim/terrarium/bots"
	"github.com/arukim/terrarium/terra"
	"log"
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
		bot := bots.Forwarder{}
		bot.Init(connectQueue)
	}

	var input string
	fmt.Scanln(&input)
}
