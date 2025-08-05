package handlers

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lenkton/51/pkg/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connectToGameUpdates(c *gin.Context) {
	game := c.MustGet("game").(*models.Game)
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	subscription_id := game.News.Subscribe(func(m models.NewsMessage) {
		data, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(1, data); err != nil {
			log.Println(err)
			return
		}
	})
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			break
		}
	}
	game.News.Unsubscribe(subscription_id)
	log.Println("seems like we have unsubscribed")
}
