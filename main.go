package main

import (
	"errors"
	"net/http"
	"strconv"

	"slices"

	"github.com/gin-gonic/gin"
)

type Game struct {
	ID            int       `json:"id"`
	CurrentPlayer *Player   `json:"currentPlayer"`
	Turns         []Turn    `json:"turns"`
	Players       []*Player `json:"players"`
}

type Turn struct {
	ID     int `json:"id"`
	dice   int
	result int
}
type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.Default()

	r.GET("/games", indexGames)
	r.POST("/games", createGame)
	r.GET("/games/:id", getGame)
	r.POST("/games/:id/join", joinGame)

	r.GET("/players", indexPlayers)

	r.Run()
}

var games = []*Game{
	{ID: 1, CurrentPlayer: nil, Players: make([]*Player, 0), Turns: make([]Turn, 0)},
}
var newGameID = 2
var players = []*Player{
	{ID: 1, Name: "Alice"},
	{ID: 2, Name: "Bob"},
}

func indexPlayers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, players)
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
func findPlayer(id int) (*Player, error) {
	for _, g := range players {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, errors.New("Game Not Found")
}
