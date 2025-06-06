// models/models.go
package models

import "gorm.io/gorm"

type TeamModel struct {
	gorm.Model
	Name         string
	Strength     int
	Points       int
	GoalsScored  int
	GoalsAgainst int
	GamesPlayed  int
	Wins         int
	Draws        int
	Losses       int
}


type MatchModel struct {
	gorm.Model
	Week      int
	Home      string
	Away      string
	HomeGoals int
	AwayGoals int
}
