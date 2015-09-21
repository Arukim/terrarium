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

type Point struct {
	X int
	Y int
}

type PlayerInfo struct {
	Player *Player
	Score  int
	Cells  map[Point]int
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
	playersInfo  []PlayerInfo
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

func (g *Game) Start() chan *Player {
	connectQueue := make(chan *Player)
	go func() {
		log.Printf("Waiting for players\n")
		g.playersInfo = make([]PlayerInfo, g.maxPlayers, g.maxPlayers)
		g.playersCount = 0
		for g.playersCount < g.maxPlayers {
			player := <-connectQueue
			log.Printf("Player %s connected\n", player.Name)
			g.playersInfo[g.playersCount].Player = player
			g.playersInfo[g.playersCount].Cells = make(map[Point]int)
			g.playersCount++
		}

		// start positions

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
	for _, playerInfo := range g.playersInfo {
		if playerInfo.Score > maxScore {
			maxScore = playerInfo.Score
		}
	}

	for _, playerInfo := range g.playersInfo {
		if playerInfo.Score == maxScore {
			log.Printf("%s is winner with %v score",
				playerInfo.Player.Name, playerInfo.Score)
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
	for i, _ := range g.playersInfo {
		g.SpawnFood()

		pInfo := &g.playersInfo[i]

		// check for old turns
		select {
		case <-pInfo.Player.PlayerTurnCh:
		default:
		}

		// start timeout, send turn summary
		var timeout = helpers.NewTimeout(g.turnTimeout)
		pInfo.Player.TurnSummaryCh <- TurnSummary{Turn: g.turn}

		// wait for player answer
		select {
		case <-pInfo.Player.PlayerTurnCh:
			log.Printf("turn was made")
			pInfo.Score++
		case <-timeout.Alarm:
			log.Printf("turn timeout")
		}
	}
	g.turn++
}
