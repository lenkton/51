package main

import (
	"lenkton/51/handlers"

	"github.com/gin-gonic/gin"
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
