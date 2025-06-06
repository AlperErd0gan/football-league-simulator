package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"github.com/AlperErd0gan/football-league-simulator/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("league.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Auto-migrate your tables
	DB.AutoMigrate(&models.TeamModel{}, &models.MatchModel{})
		// 🔥 Wipe existing data (reset DB on startup)
	DB.Exec("DELETE FROM team_models")
	DB.Exec("DELETE FROM match_models")


	// If I want to reset the auto-increment counters
	 DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'team_models'")
	 DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'match_models'")

}

