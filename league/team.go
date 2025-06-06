package league

type Team struct {
	Name         string
	Strength     int
	Points       int
	GoalsScored  int
	GoalsAgainst int
	GamesPlayed  int
	Wins         int
	Draws        int
	Losses       int
}


func (t *Team) GoalDifference() int {
	return t.GoalsScored - t.GoalsAgainst
}
