package league

import (
	"fmt"
	"sort"
)

type League struct {
	Teams     []*Team
	Simulator Simulator
	Week      int
	Fixtures  [][]Match
	Results   []Match
}

func NewLeague(teams []*Team, simulator Simulator) *League {
	return &League{
		Teams:     teams,
		Simulator: simulator,
		Week:      0,
		Fixtures:  generateFixtures(teams),
	}
}

func generateFixtures(teams []*Team) [][]Match {
	n := len(teams)
	copyTeams := make([]*Team, n)
	copy(copyTeams, teams)

	var fixtures [][]Match

	for week := 0; week < n-1; week++ {
		var weekMatches []Match
		for i := 0; i < n/2; i++ {
			home := copyTeams[i]
			away := copyTeams[n-1-i]
			weekMatches = append(weekMatches, Match{Week: week + 1, Home: home, Away: away})
		}
		// rotate
		temp := copyTeams[1]
		for i := 1; i < n-1; i++ {
			copyTeams[i] = copyTeams[i+1]
		}
		copyTeams[n-1] = temp

		fixtures = append(fixtures, weekMatches)
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
		l.Results = append(l.Results, result)
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
		if l.Week >= len(l.Fixtures) {
			fmt.Println("All weeks completed.")
		}

		matches := l.Fixtures[l.Week]
		for _, m := range matches {
			result := l.Simulator.SimulateMatch(m.Home, m.Away)
			l.Results = append(l.Results, result)
		}
		l.Week++
	}
	fmt.Println("\n✅ All weeks completed.")
	fmt.Println("\n🏁 Final League Standings:")
	l.PrintTable()
}
