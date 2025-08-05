package turn

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/lenkton/51/pkg/models"
)

func MakeTurn(game *models.Game, player *models.Player, dice int) (*models.Turn, error) {
	if player != game.CurrentPlayer {
		return nil, errors.New("it is another player's turn")
	}

	turn := &models.Turn{
		Dice:   dice,
		Result: rand.Intn(dice) + 1,
	}

	turn, err := models.MainStorage.CreateTurn(turn, game)
	if err != nil {
		return nil, fmt.Errorf("MakeTurn: error saving turn: %v", err)
	}

	// TODO: move this somewhere or call some fake save on this
	nextPlayerIndex := len(game.Turns) % len(game.Players)
	game.CurrentPlayer = game.Players[nextPlayerIndex]

	game.News.Publish(models.NewsMessage{
		"type": "newTurn",
		"turn": turn,
	})

	return turn, err
}
