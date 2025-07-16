package handlers

import (
	"lenkton/51/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindPlayersAPI(r *gin.Engine) {
	r.GET("/players", indexPlayers)
}

func indexPlayers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.MainStorage.AllPlayers())
}
