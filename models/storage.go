package models

import (
	"errors"
	"math/rand"
	"strconv"
)

type Storage struct {
	players      []*Player
	lastPlayerID int
	games        []*Game
	lastGameID   int
	lastTurnID   int
}

var MainStorage = Storage{
	players: []*Player{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	},
	lastPlayerID: 2,
	games: []*Game{
		{
			ID:            1,
			CurrentPlayer: nil,
			Players:       make([]*Player, 0),
			Turns:         make([]*Turn, 0),
			News:          &NewsCenter{callbacks: make(map[int]func(NewsMessage))},
			Status:        Created,
		},
	},
	lastGameID: 1,
	lastTurnID: 0,
}

func (s *Storage) CreateGame() *Game {
	game := Game{
		Turns:   make([]*Turn, 0),
		Players: make([]*Player, 0),
		ID:      s.lastGameID + 1,
		News:    NewNewsCenter(),
		Status:  Created,
	}

	s.lastGameID++
	s.games = append(s.games, &game)

	return &game
}

func (s *Storage) AllGames() []*Game {
	return s.games
}

func (s *Storage) FindGame(id string) (*Game, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	for _, g := range s.games {
		if g.ID == intID {
			return g, nil
		}
	}
	return nil, errors.New("Game Not Found")
}

func (s *Storage) CreatePlayer(name string) *Player {
	player := Player{
		ID:   s.lastPlayerID + 1,
		Name: name,
	}

	s.lastPlayerID++
	s.players = append(s.players, &player)

	return &player
}

func (s *Storage) AllPlayers() []*Player {
	return s.players
}

func (s *Storage) FindPlayer(id int) (*Player, error) {
	for _, g := range s.players {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, errors.New("Player Not Found")
}

// should i check the game pointer?
func (s *Storage) CreateTurn(game *Game, player *Player, dice int) (*Turn, error) {
	if player != game.CurrentPlayer {
		return nil, errors.New("it is another player's turn")
	}

	turn := Turn{
		ID:     s.lastTurnID + 1,
		Dice:   dice,
		Result: rand.Intn(dice) + 1,
	}

	s.lastTurnID++
	game.Turns = append(game.Turns, &turn)

	nextPlayerIndex := len(game.Turns) % len(game.Players)
	game.CurrentPlayer = game.Players[nextPlayerIndex]

	game.News.Publish(NewsMessage{
		"type": "newTurn",
		"turn": turn,
	})

	return &turn, nil
}
