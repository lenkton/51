package models

import "math/rand"

// TODO: make Dice and Result private
type Turn struct {
	ID     int `json:"id"`
	Dice   int
	Result int
}

func CreateTurn(game *Game, dice int) *Turn {
	turn := Turn{
		ID:     newTurnID,
		Dice:   dice,
		Result: rand.Intn(dice) + 1,
	}

	newTurnID++
	game.Turns = append(game.Turns, &turn)

	return &turn
}
