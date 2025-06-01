package main

import (
	"encoding/json"
	"net/http"
	"football-league/league"
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
	json.NewEncoder(w).Encode(leagueInstance.Teams)
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