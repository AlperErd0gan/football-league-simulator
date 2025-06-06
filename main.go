package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/AlperErd0gan/football-league-simulator/league"
	"github.com/AlperErd0gan/football-league-simulator/models"
)

func main() {
	InitDB()

	var teamModels []models.TeamModel
	DB.Find(&teamModels)

	if len(teamModels) == 0 {
		initial := []models.TeamModel{
			{Name: "Arsenal", Strength: 90},
			{Name: "Leeds", Strength: 70},
			{Name: "Everton", Strength: 80},
			{Name: "Luton", Strength: 50},
		}
		for _, t := range initial {
			DB.Create(&t)
		}
		teamModels = initial
	}

	var teams []*league.Team
	for _, t := range teamModels {
		teams = append(teams, &league.Team{
			Name:         t.Name,
			Strength:     t.Strength,
			Points:       t.Points,
			GoalsScored:  t.GoalsScored,
			GoalsAgainst: t.GoalsAgainst,
		})
	}

	l := league.NewLeague(teams, &league.StrengthBasedSimulator{}, DB)

	initAPI(l)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// ✅ Use PORT from environment for Render compatibility
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	fmt.Println("🚀 Server running on port:", port)
	http.ListenAndServe(":"+port, nil)
}
