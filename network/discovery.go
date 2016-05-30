package network

import (
	"github.com/arukim/terrarium/helpers"
	"github.com/bndr/gopencils"
	"log"
	"time"
)

type DiscoveryService struct {
	TimeoutSec int
}

func (d *DiscoveryService) Start() {
	log.Print("[discovery] service started")
	go func() {
		for {
			d.discover()
		}
	}()
}

func (d *DiscoveryService) discover() {
	log.Print("[discovery] heartbeat")
	serverFound := make(chan bool)
	var timeout = helpers.NewTimeout(time.Duration(d.TimeoutSec) * time.Second)
	go discoverServer(serverFound)
	log.Print("[discovery] w8 for response")
	select {
	case <-serverFound:
		log.Print("[discovery] server found")
		time.Sleep(10 * time.Second)
	case <-timeout.Alarm:
		log.Print("[discovery] timeout")
	}
}

type discoverRequest struct {
	Uri         string
	Type        string
	Description string
	Title       string
}

type discoverRespose struct {
	Message string
}

func discoverServer(serverFound chan bool) {
	api := gopencils.Api("http://localhost:5000/api")

	req := &discoverRequest{
		Uri:         "http://localhost:3000",
		Type:        "Terrarium.Go",
		Title:       "Terrarium Game",
		Description: "MMO for bots",
	}

	resp := &discoverRespose{}
	var res, err = api.Res("gameservers", resp).Put(req)
	if err == nil {
		log.Print(res.Response)
		serverFound <- true
	} else {
		log.Print(err)
	}
}
