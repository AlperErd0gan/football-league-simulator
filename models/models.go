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
}

type MatchModel struct {
	gorm.Model
	Week      int
	Home      string
	Away      string
	HomeGoals int
	AwayGoals int
}
