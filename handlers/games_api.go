package handlers

import (
	"fmt"
	"lenkton/51/models"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
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
		withGame.POST("/start", startGame)
	}
}

func indexGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.AllGames())
}

func joinGame(c *gin.Context) {
	game := c.MustGet("game").(*models.Game)
	if !game.JoinAllowed() {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "joining to this game is not allowed"},
		)
		return
	}
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
	if !game.CanMakeTurns() {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot make a turn now"})
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
	turn := game.CreateTurn(requestBody.Dice)
	game.News.Publish(models.NewsMessage{
		"type": "newTurn",
		"turn": turn,
	})
	c.IndentedJSON(http.StatusOK, turn)
}

func startGame(c *gin.Context) {
	game := c.MustGet("game").(*models.Game)
	err := game.Start()
	if err == nil {
		c.JSON(http.StatusOK, game)
	} else {
		// TODO: sane error messages
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}
}
