package bots

import (
	. "github.com/arukim/terrarium/terra"
	"log"
	"math/rand"
	"time"
)

var _possibleNames = [...]string{"Ivan", "Drake", "Sussana", "NyanCat", "David Blain", "Sub-Zero", "Yakobovich"}

type TerraBot interface {
	Init(connectQueue chan *Player)
}

type Forwarder struct {
	ConnectQueue chan *Player
}

func (f Forwarder) Init(connectQueue chan *Player) {
	f.ConnectQueue = connectQueue
	player := NewPlayer(_possibleNames[rand.Intn(len(_possibleNames))])

	connectQueue <- player
	go func() {
		//player logic
		for {
			tInfo := <-player.TurnSummaryCh
			log.Printf("%s: got %d turn info\n", player.Name, tInfo.Turn)
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			player.PlayerTurnCh <- PlayerTurn{}
		}
	}()
}
