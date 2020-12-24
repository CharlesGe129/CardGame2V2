package main

import (
	"github.com/CharlesGe129/CardGame2V2/pkg"
	"github.com/CharlesGe129/CardGame2V2/pkg/player"
)

func main() {
	players := [4]player.Player{
		player.NewPlayer("玩家1", 1),
		player.NewPlayer("玩家2", 2),
		player.NewPlayer("玩家3", 1),
		player.NewPlayer("玩家4", 2),
	}
	game := pkg.NewGame(players)
	game.Start()
}
