package main

import (
	"fmt"
	"net/http"
	"github.com/AlperErd0gan/football-league-simulator/league"
    "github.com/AlperErd0gan/football-league-simulator/models"
)

func main() {

	InitDB()

	var teamModels []models.TeamModel
	DB.Find(&teamModels)

	// If DB is empty, insert initial teams
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

	// Convert to []*league.Team
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
	fmt.Println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
