package main

import (
	"encoding/json"
	"net/http"
	"football-league/league"
	"sort"
	"football-league/models" 
)


var leagueInstance *league.League

func initAPI(l *league.League) {
	leagueInstance = l

	http.HandleFunc("/league", getLeague)
	http.HandleFunc("/play/week", playNextWeek)
	http.HandleFunc("/play/all", playAllWeeks)
	http.HandleFunc("/week", getCurrentWeek) // initAPI
	http.HandleFunc("/restart", restartLeague)



	http.HandleFunc("/debug/db", func(w http.ResponseWriter, r *http.Request) {
		var teams []models.TeamModel
		var matches []models.MatchModel
		DB.Find(&teams)
		DB.Find(&matches)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"teams":   teams,
			"matches": matches,
		})
	})
}

func getLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Sort the table
	teams := leagueInstance.Teams
	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points == teams[j].Points {
			return teams[i].GoalDifference() > teams[j].GoalDifference()
		}
		return teams[i].Points > teams[j].Points
	})

	// Find the last played week from the Results slice
	lastWeek := 0
	for _, match := range leagueInstance.Results {
		if match.Week > lastWeek {
			lastWeek = match.Week
		}
	}


	var weekResults []league.Match
	for _, match := range leagueInstance.Results {
		if match.Week == lastWeek {
			weekResults = append(weekResults, match)
		}
	}

	// Predictions (simple % based on points)
	totalPoints := 0
	for _, team := range teams {
		totalPoints += team.Points
	}
	predictions := make(map[string]int)
	for _, team := range teams {
		if totalPoints > 0 {
			predictions[team.Name] = int(float64(team.Points)/float64(totalPoints)*100 + 0.5)
		} else {
			predictions[team.Name] = 25 // default
		}
	}



	// Final response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"week":         leagueInstance.Week,
		"table":        teams,
		"matchResults": weekResults,
		"predictions":  predictions,
	})
}



func playNextWeek(w http.ResponseWriter, r *http.Request) {
	played := leagueInstance.PlayNextWeek()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"week":   leagueInstance.Week,
		"played": played,
	})
}

func playAllWeeks(w http.ResponseWriter, r *http.Request) {
	leagueInstance.PlayAllWeeks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All weeks played",
	})
}


func getCurrentWeek(w http.ResponseWriter, r *http.Request) {
	var lastMatch models.MatchModel
	result := DB.Order("week desc").First(&lastMatch)

	week := 0
	if result.Error == nil {
		week = lastMatch.Week
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"week": week,
	})
}


func restartLeague(w http.ResponseWriter, r *http.Request) {
	// Reset DB

	DB.Exec("DELETE FROM team_models")
	DB.Exec("DELETE FROM match_models")


	DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'team_models'")
	DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'match_models'")


	// Re-insert original teams
	initial := []models.TeamModel{
		{Name: "Arsenal", Strength: 90},
		{Name: "Leeds", Strength: 70},
		{Name: "Everton", Strength: 80},
		{Name: "Luton", Strength: 50},
	}
	for _, t := range initial {
		DB.Create(&t)
	}

	// Reset in-memory league
	var teams []*league.Team
	for _, t := range initial {
		teams = append(teams, &league.Team{
			Name:         t.Name,
			Strength:     t.Strength,
			Points:       0,
			GoalsScored:  0,
			GoalsAgainst: 0,
		})
	}

	leagueInstance = league.NewLeague(teams, &league.StrengthBasedSimulator{}, DB)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "League restarted"})
}

