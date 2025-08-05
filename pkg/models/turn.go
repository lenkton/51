package models

// TODO: make Dice and Result private
type Turn struct {
	ID     int `json:"id"`
	Dice   int `json:"dice"`
	Result int `json:"result"`
}
