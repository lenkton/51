package models

import "errors"

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func AllPlayers() []*Player {
	return players
}

func FindPlayer(id int) (*Player, error) {
	for _, g := range players {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, errors.New("Player Not Found")
}

func CreatePlayer(name string) *Player {
	player := Player{
		ID:   newPlayerID,
		Name: name,
	}

	newPlayerID++
	players = append(players, &player)

	return &player
}
