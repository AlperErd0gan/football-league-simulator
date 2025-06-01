package main

import (
	"encoding/json"
	"net/http"
	"football-league/league"
	"sort"
)


var leagueInstance *league.League

func initAPI(l *league.League) {
	leagueInstance = l

	http.HandleFunc("/league", getLeague)
	http.HandleFunc("/play/week", playNextWeek)
	http.HandleFunc("/play/all", playAllWeeks)
}

func getLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teams := leagueInstance.Teams
	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points == teams[j].Points {
			return teams[i].GoalDifference() > teams[j].GoalDifference()
		}
		return teams[i].Points > teams[j].Points
	})

	// Last played week's results
	lastWeek := leagueInstance.Week
	if lastWeek > 0 {
		lastWeek--
	}
	var weekResults []league.Match
	for _, match := range leagueInstance.Results {
		if match.Week == lastWeek {
			weekResults = append(weekResults, match)
		}
	}

	// Prediction
	totalPoints := 0
	for _, team := range teams {
		totalPoints += team.Points
	}
	predictions := make(map[string]int)
	for _, team := range teams {
		if totalPoints > 0 {
			predictions[team.Name] = int(float64(team.Points)/float64(totalPoints)*100 + 0.5)
		} else {
			predictions[team.Name] = 25 // even chance before games
		}
	}

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