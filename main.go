package main

import (
	"fmt"
	"net/http"
	"football-league/league"
)

func main() {
	teams := []*league.Team{
		{Name: "Arsenal", Strength: 90},
		{Name: "Leeds", Strength: 70},
		{Name: "Everton", Strength: 80},
		{Name: "Luton", Strength: 50},
	}
	l := league.NewLeague(teams, &league.StrengthBasedSimulator{})

	initAPI(l)

	fmt.Println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
