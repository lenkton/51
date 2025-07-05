package handlers

import (
	"encoding/json"
	"fmt"
	"lenkton/51/models"
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func BindGamesAPI(r *gin.Engine) {
	r.GET("/games", indexGames)
	r.POST("/games", createGame)
	withGame := r.Group("/games/:id", middlewareFindGame)
	{
		withGame.GET("", getGame)
		withGame.POST("/join", joinGame)
		withGame.POST("/roll", rollDice)
		withGame.GET("/updates", connectToGameUpdates)
	}
}

func indexGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.AllGames())
}

func joinGame(c *gin.Context) {
	game := c.MustGet("game").(*models.Game)
	player, err := findUser(c)
	if err == nil {
		joinAuthed(c, game, player)
	} else {
		joinUnauthed(c, game)
	}
}

func joinAuthed(c *gin.Context, game *models.Game, player *models.Player) {
	if !slices.Contains(game.Players, player) {
		game.Players = append(game.Players, player)
	}
	c.IndentedJSON(http.StatusOK, game)
}

func joinUnauthed(c *gin.Context, game *models.Game) {
	var requestBody struct {
		UserName string `json:"userName"`
	}
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
	game := c.MustGet("game").(*models.Game)
	user, _ := findUser(c)
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

func createGame(c *gin.Context) {
	game := models.CreateGame()

	c.IndentedJSON(http.StatusCreated, game)
}

func middlewareFindGame(c *gin.Context) {
	id := c.Param("id")
	game, err := models.FindGame(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	c.Set("game", game)
}

func findUser(c *gin.Context) (*models.Player, error) {
	userStringID, err := c.Cookie("user_id")
	var user *models.Player
	if err == nil {
		var userID int
		fmt.Sscan(userStringID, &userID)
		user, err = models.FindPlayer(userID)
	}
	return user, err
}

func rollDice(c *gin.Context) {
	game := c.MustGet("game").(*models.Game)
	var requestBody struct {
		Dice int `json:"dice"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// TODO: add sane error messages
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	turn := models.CreateTurn(game, requestBody.Dice)
	game.News.Publish(models.NewsMessage{
		"type": "newTurn",
		"turn": turn,
	})
	c.IndentedJSON(http.StatusOK, turn)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connectToGameUpdates(c *gin.Context) {
	game := c.MustGet("game").(*models.Game)
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
