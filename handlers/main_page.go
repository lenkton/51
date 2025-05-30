package handlers

import (
	"lenkton/51/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindMainPage(r *gin.Engine) {
	r.GET("/", getMainPage)
}

func getMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "main_page/index.html", gin.H{
		"games": models.AllGames(),
	})
}
