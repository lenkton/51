package models

// TODO: make Dice and Result private
type Turn struct {
	ID     int `json:"id"`
	Dice   int
	Result int
}
