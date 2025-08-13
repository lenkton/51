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
	ID            int             `json:"id"`
	CurrentPlayer *Player         `json:"currentPlayer"`
	Winner        *Player         `json:"winner,omitempty"`
	Turns         map[int][]*Turn `json:"turns"`
	Players       []*Player       `json:"players"`
	ActivePlayers []*Player       `json:"active_players"`
	News          *NewsCenter     `json:"-"`
	Status        gameStatus      `json:"status"`
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
	game.ActivePlayers = append(game.ActivePlayers, player)

	game.News.Publish(NewsMessage{
		"type":   "newPlayer",
		"player": player,
	})
}

func (game *Game) CanMakeTurns(p *Player) bool {
	if game.Status == Started && game.CurrentPlayer.ID == p.ID {
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
	for _, turn := range game.Turns[player.ID] {
		res += turn.Result
	}
	return res
}

func (game *Game) CompleteWithWinner(player *Player) {
	game.ActivePlayers = []*Player{}
	game.Status = Finished
	game.Winner = player
}

func (game *Game) CompleteWithNoWinner() {
	game.ActivePlayers = []*Player{}
	game.Status = Finished
}

// WARN: it assumes, that the previous active player is still in the game
func (game *Game) MoveToNextPlayer() {
	previousCurrentPlayer := game.CurrentPlayer
	previousIndex := slices.Index(game.ActivePlayers, previousCurrentPlayer)
	// TODO: it should be an error, I suppose...
	if previousIndex < 0 {
		panic("Error: MoveToNextPlayer: current player is not an active player")
	}
	nextIndex := (previousIndex + 1) % len(game.ActivePlayers)

	game.CurrentPlayer = game.ActivePlayers[nextIndex]
}
