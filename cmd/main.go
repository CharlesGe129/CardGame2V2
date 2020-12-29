package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/CharlesGe129/CardGame2V2/pkg"
	"github.com/CharlesGe129/CardGame2V2/pkg/player"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLog() {
	var out = zerolog.NewConsoleWriter()
	out.TimeFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zlog.Logger = zlog.Output(out).With().Caller().Logger()

	// set log level
	level, err := strconv.ParseUint(os.Getenv("LOG_LEVEL"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	zerolog.SetGlobalLevel(zerolog.Level(level))
}

func main() {
	InitLog()

	players := [4]player.Player{
		player.NewHumanPlayer("玩家1", 1),
		player.NewAiPlayer("电脑1", 2),
		player.NewAiPlayer("电脑2", 1),
		player.NewAiPlayer("电脑3", 2),
	}
	game := pkg.NewGame(players, "K")
	game.Start()
}
