package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lenkton/51/pkg/models"
)

func BindPlayersAPI(r *gin.Engine) {
	r.GET("/players", indexPlayers)
}

func indexPlayers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.MainStorage.AllPlayers())
}
