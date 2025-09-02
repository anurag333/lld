package engine

import "awesomeProject/snake_n_ladder/model"

func InitEngine(numberOfSnakes, numberOfLadders, boardSize int64) *engine {
	return &engine{
		board:   model.InitBoard(boardSize, numberOfSnakes, numberOfLadders),
		dice:    model.InitDice(6),
		players: []*model.Player{},
	}
}
