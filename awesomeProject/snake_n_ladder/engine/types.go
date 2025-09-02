package engine

import "awesomeProject/snake_n_ladder/model"

type engine struct {
	board   *model.Board
	dice    *model.Dice
	players []*model.Player
}
