package models

var games = []*Game{
	{ID: 1, CurrentPlayer: nil, Players: make([]*Player, 0), Turns: make([]Turn, 0)},
}

var newGameID = 2
var newPlayerID = 3
var players = []*Player{
	{ID: 1, Name: "Alice"},
	{ID: 2, Name: "Bob"},
}
