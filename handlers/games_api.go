package handlers

import (
	"fmt"
	"lenkton/51/models"
	"net/http"

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
	// UserID int `json:"userId"`
	UserName string `json:"userName"`
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
		// TODO: add sane error messages
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	// TODO: check if the player has already entered the game
	player := models.CreatePlayer(requestBody.UserName)
	game.Players = append(game.Players, player)
	c.SetCookie("user_id", fmt.Sprint(player.ID), 1000000, "/", "localhost", false, true)
	c.SetCookie("user_name", player.Name, 1000000, "/", "localhost", false, true)

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
	userStringID, err := c.Cookie("user_id")
	var user *models.Player
	if err == nil {
		var userID int
		fmt.Sscan(userStringID, &userID)
		user, err = models.FindPlayer(userID)
		if err != nil {
			user = nil
		}
	}
	turns := game.Turns
	rounds := make([][]*models.Turn, 0)
	for i, turn := range turns {
		if i%len(game.Players) == 0 {
			rounds = append(rounds, make([]*models.Turn, 0))
		}
		rounds[len(rounds)-1] = append(rounds[len(rounds)-1], turn)
	}
	c.HTML(http.StatusOK, "games/show.html", gin.H{
		"game":   game,
		"user":   user,
		"rounds": rounds,
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
