package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	bindGamesAPI(r)
	bindPlayersAPI(r)

	r.Run()
}
