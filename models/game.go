package models

import (
	"errors"
	"fmt"
	"strconv"
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

func AllGames() []*Game {
	return games
}

func FindGame(id string) (*Game, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	for _, g := range games {
		if g.ID == intID {
			return g, nil
		}
	}
	return nil, errors.New("Game Not Found")
}

func CreateGame() *Game {
	game := Game{
		Turns:   make([]*Turn, 0),
		Players: make([]*Player, 0),
		ID:      newGameID,
		News:    NewNewsCenter(),
		Status:  Created,
	}

	newGameID++
	games = append(games, &game)

	return &game
}

func (game *Game) Start() error {
	if game.Status != Created {
		return fmt.Errorf("cannot start a game in the %s state", game.Status)
	}

	game.Status = Started
	return nil
}

func (game *Game) JoinAllowed() bool {
	if game.Status == Created {
		return true
	} else {
		return false
	}
}

func (game *Game) CanMakeTurns() bool {
	if game.Status == Started {
		return true
	} else {
		return false
	}
}
