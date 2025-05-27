package main

import (
	"errors"
	"net/http"

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

	bindGamesAPI(r)

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

func findPlayer(id int) (*Player, error) {
	for _, g := range players {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, errors.New("Game Not Found")
}
