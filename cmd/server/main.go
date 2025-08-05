package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lenkton/51/pkg/handlers"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("views/**/*")

	r.Static("/static", "./static")

	handlers.BindMainPage(r)
	handlers.BindGamesAPI(r)
	handlers.BindPlayersAPI(r)

	r.Run()
}
