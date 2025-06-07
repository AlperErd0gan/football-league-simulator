package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"github.com/AlperErd0gan/football-league-simulator/league"
    "github.com/AlperErd0gan/football-league-simulator/models"
	"errors"
	"gorm.io/gorm" 
	"log"
)



var leagueInstance *league.League

func initAPI(l *league.League) {
	leagueInstance = l

	http.HandleFunc("/league", getLeague)
	http.HandleFunc("/play/week", playNextWeek)
	http.HandleFunc("/play/all", playAllWeeks)
	http.HandleFunc("/week", getCurrentWeek) // initAPI
	http.HandleFunc("/restart", restartLeague)
	http.HandleFunc("/results/all", getAllMatchResults)

	http.HandleFunc("/edit/match", editMatchResult)




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

	var teamModels []models.TeamModel
	DB.Find(&teamModels)

	var teams []*league.Team
	for _, t := range teamModels {
		teams = append(teams, &league.Team{
			Name:         t.Name,
			Strength:     t.Strength,
			Points:       t.Points,
			GoalsScored:  t.GoalsScored,
			GoalsAgainst: t.GoalsAgainst,
			GamesPlayed:  t.GamesPlayed,
			Wins:         t.Wins,
			Draws:        t.Draws,
			Losses:       t.Losses,
			

		})
	}
	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points == teams[j].Points {
			return teams[i].GoalDifference() > teams[j].GoalDifference()
		}
		return teams[i].Points > teams[j].Points
	})


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

	// Predictions we get simple % based on points
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
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Only log unexpected errors
		log.Println("Unexpected DB error:", result.Error)
	}
	if result.Error == nil {
		week = lastMatch.Week
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"week": week,
	})
}



func restartLeague(w http.ResponseWriter, r *http.Request) {
	// Reset DB to restart the league

	DB.Exec("DELETE FROM team_models")
	DB.Exec("DELETE FROM match_models")


	DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'team_models'")
	DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'match_models'")



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

func getAllMatchResults(w http.ResponseWriter, r *http.Request) {
	var matches []models.MatchModel
	DB.Order("week asc").Find(&matches)

	grouped := make(map[int][]models.MatchModel)
	for _, m := range matches {
		grouped[m.Week] = append(grouped[m.Week], m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grouped)
}

func editMatchResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.MatchModel
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var match models.MatchModel
	if err := DB.First(&match, input.ID).Error; err != nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	// Update match with new values
	match.HomeGoals = input.HomeGoals
	match.AwayGoals = input.AwayGoals
	DB.Save(&match)

	// Recalculate standings from scratch
	recalculateLeague()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "match updated and standings recalculated"})
}


func recalculateLeague() {
	var teams []models.TeamModel
	DB.Find(&teams)
	for _, team := range teams {
		team.Points = 0
		team.GoalsScored = 0
		team.GoalsAgainst = 0
		DB.Save(&team)
	}

	var matches []models.MatchModel
	DB.Order("week asc").Find(&matches)
	for _, match := range matches {
		updateTeamStats(match)
	}

	// Reset  league to reflect updated DB
	var newTeams []*league.Team
	for _, t := range teams {
		newTeams = append(newTeams, &league.Team{
			Name:         t.Name,
			Strength:     t.Strength,
			Points:       t.Points,
			GoalsScored:  t.GoalsScored,
			GoalsAgainst: t.GoalsAgainst,
		})
	}
	leagueInstance.Teams = newTeams
}

func updateTeamStats(match models.MatchModel) {
	var home, away models.TeamModel
	DB.Where("name = ?", match.Home).First(&home)
	DB.Where("name = ?", match.Away).First(&away)

	// Update Points
	if match.HomeGoals > match.AwayGoals {
		home.Points += 3
	} else if match.AwayGoals > match.HomeGoals {
		away.Points += 3
	} else {
		home.Points += 1
		away.Points += 1
	}

	// Update Goals
	home.GoalsScored += match.HomeGoals
	home.GoalsAgainst += match.AwayGoals
	away.GoalsScored += match.AwayGoals
	away.GoalsAgainst += match.HomeGoals

	// Save Updates
	DB.Save(&home)
	DB.Save(&away)
}


