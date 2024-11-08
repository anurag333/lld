package snakeandladder

type Snake struct {
	Start int
	End   int
}

type Ladder struct {
	Start int
	End   int
}

type Player struct {
	ID   string
	Name string
}

type SnakeAndLadderBoard struct {
	Size         int
	Snakes       []Snake
	Ladders      []Ladder
	PlayerPieces map[string]int
}

func NewSnakeAndLadderBoard(size int) *SnakeAndLadderBoard {
	return &SnakeAndLadderBoard{
		Size:         size,
		PlayerPieces: make(map[string]int),
	}
}
