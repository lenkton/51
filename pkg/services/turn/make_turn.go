package turn

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"

	"github.com/lenkton/51/pkg/models"
)

func MakeTurn(game *models.Game, player *models.Player, dice int) (*models.Turn, error) {
	if player != game.CurrentPlayer {
		return nil, errors.New("it is another player's turn")
	}

	turn := buildTurn(dice)

	turn, err := models.MainStorage.CreateTurn(turn, game)
	if err != nil {
		return nil, fmt.Errorf("MakeTurn: error saving turn: %v", err)
	}

	checkPlayerTotal(player, game)

	updateNextPlayer(game)
	publishNewTurn(turn, game)

	return turn, err
}

func buildTurn(dice int) *models.Turn {
	return &models.Turn{
		Dice:   dice,
		Result: rand.Intn(dice) + 1,
	}
}

const WinNumber = 51

// WARN: MustPlayerTotal will panic if the player is not in the game
func checkPlayerTotal(player *models.Player, game *models.Game) {
	total := game.MustPlayerTotal(player)
	if total == WinNumber {
		game.CompleteWithWinner(player)
		return
	}
	if total > WinNumber {
		game.ActivePlayers = slices.DeleteFunc(
			game.ActivePlayers, func(p *models.Player) bool {
				return p == player
			},
		)
	}
	if len(game.ActivePlayers) == 0 {
		game.CompleteWithNoWinnter()
	}
}

func updateNextPlayer(game *models.Game) {
	// TODO: move this somewhere or call some fake save on this
	// TODO: count only active players
	nextPlayerIndex := len(game.Turns) % len(game.Players)
	game.CurrentPlayer = game.Players[nextPlayerIndex]
}

func publishNewTurn(turn *models.Turn, game *models.Game) {
	game.News.Publish(models.NewsMessage{
		"type": "newTurn",
		"turn": turn,
	})
}
