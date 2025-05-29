package models

type Turn struct {
	ID     int `json:"id"`
	dice   int
	result int
}
