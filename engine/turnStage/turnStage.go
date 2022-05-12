package turnStage

type TurnStage int

const (
	Draw        TurnStage = 0
	PlaceTile   TurnStage = 1
	PlaceMeeple TurnStage = 2
	Score       TurnStage = 3
	Pass        TurnStage = 4
)
