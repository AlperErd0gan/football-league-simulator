package league

import (
	"fmt"
	"sort"
	"gorm.io/gorm"
	"github.com/AlperErd0gan/football-league-simulator/models"
)

type League struct {
	Teams     []*Team
	Simulator Simulator
	Week      int
	Fixtures  [][]Match
	Results   []Match
	DB        *gorm.DB 
}

func NewLeague(teams []*Team, simulator Simulator, db *gorm.DB) *League {
	return &League{
		Teams:     teams,
		Simulator: simulator,
		Week:      0,
		Fixtures:  generateFixtures(teams),
		DB:        db,
	}
}

func generateFixtures(teams []*Team) [][]Match {
	n := len(teams)
	var fixtures [][]Match

	// First leg 
	for week := 0; week < n-1; week++ {
		var weekMatches []Match
		for i := 0; i < n/2; i++ {
			home := teams[i]
			away := teams[n-1-i]
			if week%2 == 0 {
				weekMatches = append(weekMatches, Match{Week: week + 1, Home: home, Away: away})
			} else {
				weekMatches = append(weekMatches, Match{Week: week + 1, Home: away, Away: home})
			}
		}
		// Rotate teams except the first one
		rotated := append([]*Team{teams[0]}, append(teams[n-1:], teams[1:n-1]...)...)
		copy(teams, rotated)
		fixtures = append(fixtures, weekMatches)
	}

	// Second leg
	offset := len(fixtures)
	for _, matches := range fixtures {
		var reverseWeek []Match
		for _, match := range matches {
			reverseWeek = append(reverseWeek, Match{
				Week:      match.Week + offset,
				Home:      match.Away,
				Away:      match.Home,
			})
		}
		fixtures = append(fixtures, reverseWeek)
	}

	return fixtures
}


func (l *League) PlayNextWeek() bool {
	if l.Week >= len(l.Fixtures) {
		fmt.Println("All weeks completed.")
		return false
	}

	matches := l.Fixtures[l.Week]
	fmt.Printf("\n▶️  Week %d Results:\n", l.Week+1)
	for _, m := range matches {
		result := l.Simulator.SimulateMatch(m.Home, m.Away)
		result.Week = l.Week + 1  // 🛠️ FIX: Set the correct week
		l.Results = append(l.Results, result)

	l.DB.Create(&models.MatchModel{
	Week:      result.Week,
	Home:      result.Home.Name,
	Away:      result.Away.Name,
	HomeGoals: result.HomeGoals,
	AwayGoals: result.AwayGoals,
	})

	l.DB.Model(&models.TeamModel{}).Where("name = ?", result.Home.Name).Updates(map[string]interface{}{
		"Points":       result.Home.Points,
		"GoalsScored":  result.Home.GoalsScored,
		"GoalsAgainst": result.Home.GoalsAgainst,
		"GamesPlayed":  result.Home.GamesPlayed,
		"Wins":         result.Home.Wins,
		"Draws":        result.Home.Draws,
		"Losses":       result.Home.Losses,
	})
	l.DB.Model(&models.TeamModel{}).Where("name = ?", result.Away.Name).Updates(map[string]interface{}{
		"Points":       result.Away.Points,
		"GoalsScored":  result.Away.GoalsScored,
		"GoalsAgainst": result.Away.GoalsAgainst,
		"GamesPlayed":  result.Away.GamesPlayed,
		"Wins":         result.Away.Wins,
		"Draws":        result.Away.Draws,
		"Losses":       result.Away.Losses,
	})
		fmt.Printf("%s %d - %d %s\n", result.Home.Name, result.HomeGoals, result.AwayGoals, result.Away.Name)
	}

	l.Week++
	l.PrintTable()
	return true
}


func (l *League) PrintTable() {
	fmt.Println("\n🏆 League Table:")
	sort.Slice(l.Teams, func(i, j int) bool {
		if l.Teams[i].Points == l.Teams[j].Points {
			return l.Teams[i].GoalDifference() > l.Teams[j].GoalDifference()
		}
		return l.Teams[i].Points > l.Teams[j].Points
	})

	fmt.Printf("%-12s | Points | GD\n", "Team")
	for _, t := range l.Teams {
		fmt.Printf("%-12s | %6d | %+3d\n", t.Name, t.Points, t.GoalDifference())
	}
}

func (l *League) PlayAllWeeks() {
	if l.Week >= len(l.Fixtures) {
		fmt.Println("All weeks already completed.")
		return
	}
	for l.Week < len(l.Fixtures) {
		matches := l.Fixtures[l.Week]
		fmt.Printf("\n▶️  Week %d Results:\n", l.Week+1)

		for _, m := range matches {
			result := l.Simulator.SimulateMatch(m.Home, m.Away)
			result.Week = l.Week + 1
			l.Results = append(l.Results, result)

			// Insert into DB
			l.DB.Create(&models.MatchModel{
				Week:      result.Week,
				Home:      result.Home.Name,
				Away:      result.Away.Name,
				HomeGoals: result.HomeGoals,
				AwayGoals: result.AwayGoals,
			})

			// Update team stats in DB
			l.DB.Model(&models.TeamModel{}).Where("name = ?", result.Home.Name).Updates(map[string]interface{}{
				"Points":       result.Home.Points,
				"GoalsScored":  result.Home.GoalsScored,
				"GoalsAgainst": result.Home.GoalsAgainst,
				"GamesPlayed":  result.Home.GamesPlayed,
				"Wins":         result.Home.Wins,
				"Draws":        result.Home.Draws,
				"Losses":       result.Home.Losses,
			})
			l.DB.Model(&models.TeamModel{}).Where("name = ?", result.Away.Name).Updates(map[string]interface{}{
				"Points":       result.Away.Points,
				"GoalsScored":  result.Away.GoalsScored,
				"GoalsAgainst": result.Away.GoalsAgainst,
				"GamesPlayed":  result.Away.GamesPlayed,
				"Wins":         result.Away.Wins,
				"Draws":        result.Away.Draws,
				"Losses":       result.Away.Losses,
			})

			fmt.Printf("%s %d - %d %s\n", result.Home.Name, result.HomeGoals, result.AwayGoals, result.Away.Name)
		}

		l.Week++
	}
	// Refresh in-memory league.teams after DB updates
		var updated []models.TeamModel
		l.DB.Find(&updated)

		var refreshed []*Team
		for _, t := range updated {
			refreshed = append(refreshed, &Team{
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
		l.Teams = refreshed
	fmt.Println("\n✅ All weeks completed.")
	fmt.Println("\n🏁 Final League Standings:")
	l.PrintTable()
}

