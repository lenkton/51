package models

import (
	"errors"
	"fmt"
	"slices"
)

type gameStatus string

const (
	Created  gameStatus = "created"
	Started  gameStatus = "started"
	Finished gameStatus = "finished"
)

type Game struct {
	ID            int         `json:"id"`
	CurrentPlayer *Player     `json:"currentPlayer"`
	Turns         []*Turn     `json:"turns"`
	Players       []*Player   `json:"players"`
	News          *NewsCenter `json:"-"`
	Status        gameStatus  `json:"status"`
}

func (game *Game) Start() error {
	if game.Status != Created {
		return fmt.Errorf("cannot start a game in the %s state", game.Status)
	}
	if len(game.Players) < 1 {
		return errors.New("cannot start a game: need at least 1 player")
	}

	game.Status = Started
	game.CurrentPlayer = game.Players[0]
	game.News.Publish(NewsMessage{
		"type": "gameStarted",
		"game": game,
	})
	return nil
}

func (game *Game) JoinAllowed() bool {
	if game.Status == Created {
		return true
	} else {
		return false
	}
}

// NOTE: maybe we could save time by introducing a flag
// that the user was created just now (so 100% not in the list)
func (game *Game) MustJoin(player *Player) {
	if player == nil {
		panic("MustJoin: nil player pointer")
	}
	if slices.Contains(game.Players, player) {
		return
	}

	game.Players = append(game.Players, player)

	game.News.Publish(NewsMessage{
		"type":   "newPlayer",
		"player": player,
	})
}

func (game *Game) CanMakeTurns() bool {
	if game.Status == Started {
		return true
	} else {
		return false
	}
}

func (game *Game) MustPlayerTotal(player *Player) int {
	playerIndex := slices.Index(game.Players, player)
	if playerIndex == -1 {
		panic("PlayerTotal: player is not in the game")
	}
	res := 0
	for i := playerIndex; i < len(game.Turns); i += len(game.Players) {
		res += game.Turns[i].Result
	}
	return res
}
