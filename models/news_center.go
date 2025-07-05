package models

import "log"

type NewsCenter struct {
	callbacks map[int]func(NewsMessage)
	lastId    int
}
type NewsMessage map[string]any

func NewNewsCenter() *NewsCenter {
	return &NewsCenter{callbacks: make(map[int]func(NewsMessage))}
}

func (nc *NewsCenter) Subscribe(callback func(NewsMessage)) int {
	log.Println(nc)
	id := nc.lastId + 1
	nc.callbacks[id] = callback
	nc.lastId = id
	return id
}
func (nc *NewsCenter) Unsubscribe(id int) {
	delete(nc.callbacks, id)
}
func (nc *NewsCenter) Publish(message NewsMessage) {
	for _, callback := range nc.callbacks {
		callback(message)
	}
}
