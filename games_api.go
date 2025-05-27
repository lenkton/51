package main

import (
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

func bindGamesAPI(r *gin.Engine) {
	r.GET("/games", indexGames)
	r.POST("/games", createGame)
	r.GET("/games/:id", getGame)
	r.POST("/games/:id/join", joinGame)
}

func indexGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, games)
}

type joinGameDTO struct {
	UserID int `json:"userId"`
}

func joinGame(c *gin.Context) {
	id := c.Param("id")
	game, err := findGame(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	var requestBody joinGameDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, err)
		return
	}
	player, err := findPlayer(requestBody.UserID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Player Not Found"})
		return
	}

	// check if the player has already entered the game
	if !slices.Contains(game.Players, player) {
		game.Players = append(game.Players, player)
	}

	c.IndentedJSON(http.StatusOK, game)
}
func getGame(c *gin.Context) {
	id := c.Param("id")
	game, err := findGame(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, game)
}
func createGame(c *gin.Context) {
	game := Game{Turns: make([]Turn, 0), Players: make([]*Player, 0), ID: newGameID}

	newGameID++
	games = append(games, &game)

	c.IndentedJSON(http.StatusCreated, game)
}

func findGame(id string) (*Game, error) {
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
