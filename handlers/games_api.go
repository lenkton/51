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
	withGame := r.Group("/games/:id", middlewareFindGame)
	{
		withGame.GET("", getGame)
		withGame.POST("/join", joinGame)
		withGame.POST("/roll", middlewareFindUser, rollDice)
		withGame.GET("/updates", connectToGameUpdates)
		withGame.POST("/start", middlewareFindUser, startGame)
	}
}

func indexGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.MainStorage.AllGames())
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
	game.MustJoin(player)
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

	player := models.MainStorage.CreatePlayer(requestBody.UserName)
	game.MustJoin(player)

	// WARN: it could be empty?...
	hostname := c.Request.URL.Hostname()
	c.SetCookie("user_id", fmt.Sprint(player.ID), 1000000, "/", hostname, false, true)
	c.SetCookie("user_name", player.Name, 1000000, "/", hostname, false, true)

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
	game := models.MainStorage.CreateGame()

	c.IndentedJSON(http.StatusCreated, game)
}

func middlewareFindGame(c *gin.Context) {
	id := c.Param("id")
	game, err := models.MainStorage.FindGame(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Game Not Found"})
		return
	}
	c.Set("game", game)
}

func middlewareFindUser(c *gin.Context) {
	user, err := findUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "only logged in users can make turns"})
		return
	}
	c.Set("user", user)
}

func findUser(c *gin.Context) (*models.Player, error) {
	userStringID, err := c.Cookie("user_id")
	var user *models.Player
	if err == nil {
		var userID int
		fmt.Sscan(userStringID, &userID)
		user, err = models.MainStorage.FindPlayer(userID)
	}
	return user, err
}

func rollDice(c *gin.Context) {
	player := c.MustGet("user").(*models.Player)
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
	turn, err := models.MainStorage.CreateTurn(game, player, requestBody.Dice)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
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
