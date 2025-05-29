package main

import (
	"lenkton/51/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func bindPlayersAPI(r *gin.Engine) {
	r.GET("/players", indexPlayers)
}

func indexPlayers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.AllPlayers())
}
