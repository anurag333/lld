package snakeandladder

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

type SnakeAndLadderService struct {
	board                            *SnakeAndLadderBoard
	players                          *list.List
	initialNumberOfPlayers           int
	isGameCompleted                  bool
	noOfDices                        int
	shouldGameContinueTillLastPlayer bool
	shouldAllowMultipleDiceRollOnSix bool
}

func NewSnakeAndLadderService(boardSize int) *SnakeAndLadderService {
	board := NewSnakeAndLadderBoard(boardSize)
	return &SnakeAndLadderService{
		board:     board,
		players:   list.New(),
		noOfDices: 1,
	}
}

func (s *SnakeAndLadderService) SetNoOfDices(noOfDices int) {
	s.noOfDices = noOfDices
}

func (s *SnakeAndLadderService) SetShouldGameContinueTillLastPlayer(continueTillLast bool) {
	s.shouldGameContinueTillLastPlayer = continueTillLast
}

func (s *SnakeAndLadderService) SetShouldAllowMultipleDiceRollOnSix(allowMultipleRoll bool) {
	s.shouldAllowMultipleDiceRollOnSix = allowMultipleRoll
}

func (s *SnakeAndLadderService) SetPlayers(players []Player) {
	s.initialNumberOfPlayers = len(players)
	for _, player := range players {
		s.players.PushBack(player)
		s.board.PlayerPieces[player.ID] = 0 // Each player starts at position 0
	}
}

func (s *SnakeAndLadderService) SetSnakes(snakes []Snake) {
	s.board.Snakes = snakes
}

func (s *SnakeAndLadderService) SetLadders(ladders []Ladder) {
	s.board.Ladders = ladders
}

func (s *SnakeAndLadderService) getNewPositionAfterSnakesAndLadders(position int) int {
	for {
		prevPosition := position
		for _, snake := range s.board.Snakes {
			if snake.Start == position {
				position = snake.End
			}
		}
		for _, ladder := range s.board.Ladders {
			if ladder.Start == position {
				position = ladder.End
			}
		}
		if prevPosition == position {
			break
		}
	}
	return position
}

func (s *SnakeAndLadderService) movePlayer(player Player, diceValue int) {
	oldPosition := s.board.PlayerPieces[player.ID]
	newPosition := oldPosition + diceValue

	if newPosition > s.board.Size {
		newPosition = oldPosition
	} else {
		newPosition = s.getNewPositionAfterSnakesAndLadders(newPosition)
	}

	s.board.PlayerPieces[player.ID] = newPosition
	fmt.Printf("%s rolled a %d and moved from %d to %d\n", player.Name, diceValue, oldPosition, newPosition)
}

func (s *SnakeAndLadderService) rollDice() int {
	return rand.Intn(6) + 1
}

func (s *SnakeAndLadderService) getTotalDiceValue() int {
	total := 0
	for i := 0; i < s.noOfDices; i++ {
		diceRoll := s.rollDice()
		total += diceRoll
		if diceRoll < 6 || !s.shouldAllowMultipleDiceRollOnSix {
			break
		}
	}
	return total
}

func (s *SnakeAndLadderService) hasPlayerWon(player Player) bool {
	return s.board.PlayerPieces[player.ID] == s.board.Size
}

func (s *SnakeAndLadderService) isGameComplete() bool {
	return s.players.Len() < s.initialNumberOfPlayers
}

func (s *SnakeAndLadderService) StartGame() {
	rand.Seed(time.Now().UnixNano())

	for !s.isGameComplete() {
		currentPlayerElement := s.players.Front()
		currentPlayer := currentPlayerElement.Value.(Player)

		diceValue := s.getTotalDiceValue()
		s.movePlayer(currentPlayer, diceValue)

		if s.hasPlayerWon(currentPlayer) {
			fmt.Printf("%s wins the game\n", currentPlayer.Name)
			s.board.PlayerPieces[currentPlayer.ID] = s.board.Size
			s.players.Remove(currentPlayerElement)
		} else {
			s.players.MoveToBack(currentPlayerElement)
		}
	}
}
