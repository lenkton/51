package main

import (
	"lenkton/51/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	handlers.BindGamesAPI(r)
	handlers.BindPlayersAPI(r)

	r.Run()
}
