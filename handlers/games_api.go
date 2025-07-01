package handlers

import (
	"lenkton/51/models"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func BindGamesAPI(r *gin.Engine) {
	r.GET("/games", indexGames)
	r.POST("/games", createGame)
	r.GET("/games/:id", getGame)
	r.POST("/games/:id/join", joinGame)
}

func indexGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.AllGames())
}

type joinGameDTO struct {
	UserID int `json:"userId"`
}

func joinGame(c *gin.Context) {
	id := c.Param("id")
	game, err := models.FindGame(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	var requestBody joinGameDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, err)
		return
	}
	player, err := models.FindPlayer(requestBody.UserID)
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
	game, err := models.FindGame(id)
	if err != nil {
		// TODO: add some 404
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	c.HTML(http.StatusOK, "games/show.html", gin.H{
		"game": game,
	})
}

func getGameJSON(c *gin.Context) {
	id := c.Param("id")
	game, err := models.FindGame(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, game)
}
func createGame(c *gin.Context) {
	game := models.CreateGame()

	c.IndentedJSON(http.StatusCreated, game)
}
