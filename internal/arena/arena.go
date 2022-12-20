package arena

import (
	"log"
	"math/rand"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Player struct {
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
}

type PlayerScore struct {
	Name  string
	Score int
}

type Arena struct {
	players     []Player
	leaderboard map[string]int
	encounters  int
	runnerCmd   string
}

func NewArena(runnerCmd string, encounters int, players ...Player) *Arena {
	return &Arena{
		players:     players,
		leaderboard: make(map[string]int),
		encounters:  encounters,
		runnerCmd:   runnerCmd,
	}
}

func (a *Arena) Run() []PlayerScore {
	rand.Seed(time.Now().UnixMilli())
	a.initScores()

	winnerChan := make(chan string)
	go a.countScore(winnerChan)
	wg := sync.WaitGroup{}

	count := 0

	for i := 0; i < len(a.players)-1; i++ {
		for j := i + 1; j < len(a.players); j++ {
			log.Println(a.players[i].Name, " vs ", a.players[j].Name)
			cmdParts := strings.Fields(a.runnerCmd)
			cmdParts = append(cmdParts, a.players[i].Cmd, a.players[i].Name, a.players[j].Cmd, a.players[j].Name)

			for e := 0; e < a.encounters; e++ {
				cmdParts := append(cmdParts, strconv.FormatInt(rand.Int63(), 10))
				wg.Add(1)
				go fight(&wg, cmdParts, winnerChan)
				count++
				if count == 10 {
					wg.Wait()
					count = 0
				}
			}
		}
	}

	wg.Wait()
	close(winnerChan)

	res := make([]PlayerScore, 0, len(a.players))
	for _, p := range a.players {
		res = append(res, PlayerScore{
			Name:  p.Name,
			Score: a.leaderboard[p.Name],
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})
	return res
}

func (a *Arena) initScores() {
	for _, p := range a.players {
		a.leaderboard[p.Name] = 0
	}
}

func (a *Arena) countScore(winnerChan <-chan string) {
	for winner := range winnerChan {
		a.leaderboard[winner]++
	}
}

func fight(wg *sync.WaitGroup, cmdParts []string, winnerChan chan<- string) {
	defer wg.Done()
	ex := exec.Command(cmdParts[0], cmdParts[1:]...)

	res, err := ex.Output()
	if err != nil {
		log.Println(string(res))
		panic(err)
	}
	parts := strings.Split(string(res), "\n")
	key := parts[len(parts)-2]
	winnerChan <- key
}
