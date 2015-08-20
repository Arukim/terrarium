package terra

import (
	"encoding/json"
	"fmt"
	"github.com/arukim/terrarium/helpers"
	"log"
	"math/rand"
	"os"
	"time"
)

type TurnSummary struct {
	Turn int
}

type PlayerTurn struct {
}

type PlayerInfo struct {
	Name          string
	TurnSummaryCh chan TurnSummary
	PlayerTurnCh  chan PlayerTurn
	Score         int
}

type Point struct {
	X int
	Y int
}

type Game struct {
	// settings
	maxPlayers    int
	maxTurns      int
	mapWidth      int
	mapHeight     int
	foodSpawnRate int
	turnTimeout   time.Duration
	// current info
	players      []PlayerInfo
	playersCount int
	turn         int
	food         map[Point]int
	// log file
	logFile *os.File
}

func NewGame(maxPlayers int, maxTurns int, turnTimeout time.Duration) *Game {
	g := new(Game)
	g.maxPlayers = maxPlayers
	g.maxTurns = maxTurns
	g.turnTimeout = turnTimeout
	g.mapWidth = 100
	g.mapHeight = 100
	g.foodSpawnRate = 20

	g.food = make(map[Point]int)

	logFile, err := os.OpenFile("turns.json", os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	g.logFile = logFile
	return g
}

func (g *Game) Start() chan PlayerInfo {
	connectQueue := make(chan PlayerInfo)
	go func() {
		log.Printf("Waiting for players\n")
		g.players = make([]PlayerInfo, g.maxPlayers, g.maxPlayers)
		g.playersCount = 0
		for g.playersCount < g.maxPlayers {
			player := <-connectQueue
			log.Printf("Player %s connected\n", player.Name)
			g.players[g.playersCount] = player
			g.playersCount++
		}

		log.Printf("Starting the game\n")
		go func() {
			for g.turn < g.maxTurns {
				g.Turn()
			}
			g.PrintFoodMap()
			g.FindWinner()

			g.logFile.Close()
		}()
	}()
	return connectQueue
}

func (g *Game) FindWinner() {
	maxScore := 0
	for _, player := range g.players {
		if player.Score > maxScore {
			maxScore = player.Score
		}
	}

	for _, player := range g.players {
		if player.Score == maxScore {
			log.Printf("%s is winner with %v score", player.Name, player.Score)
		}
	}
}

func (g *Game) SpawnFood() {
	var foodDiff = make(map[Point]int)
	for i := 0; i < g.foodSpawnRate; i++ {
		var point = Point{X: rand.Intn(g.mapWidth),
			Y: rand.Intn(g.mapHeight)}
		foodDiff[point] = 1
	}
	LogFoodDiff(foodDiff, g.logFile)

	for point, value := range foodDiff {
		var currValue = g.food[point]
		currValue += value
		g.food[point] = currValue
	}
}

func LogFoodDiff(diff map[Point]int, logFile *os.File) {
	tmpDiff := make(map[string]int)
	for point, value := range diff {
		var arg = fmt.Sprintf("%v,%v", point.X, point.Y)
		tmpDiff[arg] = value
	}
	json, _ := json.Marshal(tmpDiff)
	_, err := logFile.Write(json)
	if err != nil {
		log.Fatal("Can't write to file")
	}
}

func (g *Game) PrintFoodMap() {
	for key, value := range g.food {
		log.Printf("%v - %v", value, key)
	}
}

func (g *Game) Turn() {
	log.Printf("Turn %d\n", g.turn)
	for i, _ := range g.players {
		g.SpawnFood()

		player := &g.players[i]

		// check for old turns
		select {
		case <-player.PlayerTurnCh:
		default:
		}

		// start timeout, send turn summary
		var timeout = helpers.NewTimeout(g.turnTimeout)
		player.TurnSummaryCh <- TurnSummary{Turn: g.turn}

		// wait for player answer
		select {
		case <-player.PlayerTurnCh:
			log.Printf("turn was made")
			player.Score++
		case <-timeout.Alarm:
			log.Printf("turn timeout")
		}
	}
	g.turn++
}
