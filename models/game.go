package models

import (
	"errors"
	"strconv"
)

type Game struct {
	ID            int       `json:"id"`
	CurrentPlayer *Player   `json:"currentPlayer"`
	Turns         []Turn    `json:"turns"`
	Players       []*Player `json:"players"`
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
		Turns:   make([]Turn, 0),
		Players: make([]*Player, 0),
		ID:      newGameID,
	}

	newGameID++
	games = append(games, &game)

	return &game
}
