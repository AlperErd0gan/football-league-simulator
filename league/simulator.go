package league

import (
	"math"
	"math/rand"
)


type Simulator interface {
	SimulateMatch(home, away *Team) Match
}

func simulateGoals(mean float64) int {
	// Poisson-inspired random scoring
	l := math.Exp(-mean)
	k := 0
	p := 1.0

	for p > l {
		k++
		p *= rand.Float64()
	}
	return k - 1
}


type StrengthBasedSimulator struct{}

func (s *StrengthBasedSimulator) SimulateMatch(home, away *Team) Match {

	// Normalize strength to scoring potential (e.g., 1.5–3.5 range)
	homeScoreFactor := float64(home.Strength) / float64(home.Strength+away.Strength)
	awayScoreFactor := float64(away.Strength) / float64(home.Strength+away.Strength)

	// Use base mean goals to simulate scoring chance
	homeMeanGoals := 1.5 + homeScoreFactor*2 // range ~1.5–3.5
	awayMeanGoals := 1.5 + awayScoreFactor*2

	homeGoals := simulateGoals(homeMeanGoals)
	awayGoals := simulateGoals(awayMeanGoals)

	home.GamesPlayed++
	away.GamesPlayed++

	// Points
	if homeGoals > awayGoals {
		home.Points += 3
		home.Wins++
		away.Losses++
	} else if awayGoals > homeGoals {
		away.Points += 3
		away.Wins++
		home.Losses++
	} else {
		home.Points += 1
		away.Points += 1
		home.Draws++
		away.Draws++
	}

	home.GoalsScored += homeGoals
	home.GoalsAgainst += awayGoals
	away.GoalsScored += awayGoals
	away.GoalsAgainst += homeGoals

	return Match{
		Home:      home,
		Away:      away,
		HomeGoals: homeGoals,
		AwayGoals: awayGoals,
	}
}


