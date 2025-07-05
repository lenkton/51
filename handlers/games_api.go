package handlers

import (
	"encoding/json"
	"fmt"
	"lenkton/51/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func BindGamesAPI(r *gin.Engine) {
	r.GET("/games", indexGames)
	r.POST("/games", createGame)
	r.GET("/games/:id", getGame)
	r.POST("/games/:id/join", joinGame)
	r.POST("/games/:id/roll", rollDice)
	r.GET("/games/:id/updates", connectToGameUpdates)
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
	game.News.Publish(models.NewsMessage{
		"type":   "newPlayer",
		"player": player,
	})
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

func rollDice(c *gin.Context) {
	id := c.Param("id")
	game, err := models.FindGame(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	var requestBody struct {
		Dice int `json:"dice"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// TODO: add sane error messages
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	turn := models.CreateTurn(game, requestBody.Dice)
	c.IndentedJSON(http.StatusOK, turn)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connectToGameUpdates(c *gin.Context) {
	id := c.Param("id")
	game, err := models.FindGame(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	subscription_id := game.News.Subscribe(func(m models.NewsMessage) {
		data, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(1, data); err != nil {
			log.Println(err)
			return
		}
	})
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			break
		}
	}
	game.News.Unsubscribe(subscription_id)
	log.Println("seems like we have unsubscribed")
}
