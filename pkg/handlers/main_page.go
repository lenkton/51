package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lenkton/51/pkg/models"
)

func BindMainPage(r *gin.Engine) {
	r.GET("/", getMainPage)
}

func getMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "main_page/index.html", gin.H{
		"games": models.MainStorage.AllGames(),
	})
}
