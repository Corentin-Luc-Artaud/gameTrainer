package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/corentin-luc-artaud/trainer/internal/arena"
)

func main() {
	playersDescription := flag.String("players", "", "path to players description json file")
	encounters := flag.Int("encounters", 3, "number of encounter between each player")
	runnerCmd := flag.String("runner", "", "cmd to run the match, runner should accept the following arguments <player1 cmd> <player1 name> <player2 cmd> <player2 name>,runner is expect to return on it's last line the name of the winner")
	flag.Parse()

	players := readPlayers(*playersDescription)

	gauntlet := arena.NewArena(*runnerCmd, *encounters, players...)
	res := gauntlet.Run()

	for i, p := range res {
		fmt.Printf("%d \t %s \t%d\n", i+1, p.Name, p.Score)
	}

}

func readPlayers(path string) []arena.Player {
	var players []arena.Player

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(file).Decode(&players)
	if err != nil {
		panic(err)
	}

	return players
}
