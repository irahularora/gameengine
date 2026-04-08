package engine

import (
	"fmt"
	"game-engine/internal/models"
	"sync"
	"time"
)

type GameEngine struct {
	winnerFound    bool
	totalSubmitted int
	totalCorrect   int
	totalIncorrect int
	startTime      time.Time
	discoveryTime  time.Duration
	mu             sync.Mutex
}

func NewGameEngine() *GameEngine {
	return &GameEngine{
		startTime: time.Now(),
	}
}

// Evaluate checks if a response is correct, updates metrics, and declares a winner if none exists.
// It also records the time taken from startup to find the first winner.
func (e *GameEngine) Evaluate(resp models.UserResponse) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.totalSubmitted++

	isCorrect := resp.Answer == models.AnswerYes
	if isCorrect {
		e.totalCorrect++
	} else {
		e.totalIncorrect++
	}

	// Ignore all subsequent responses once a winner is found
	if !e.winnerFound && isCorrect {
		e.winnerFound = true
		e.discoveryTime = time.Since(e.startTime)
		fmt.Printf("\nWINNER FOUND! User ID: %s has won the game! (Time: %v)\n", resp.UserID, e.discoveryTime)
	}

	return isCorrect
}

// GetMetrics returns the current submission statistics and discovery time safely.
func (e *GameEngine) GetMetrics() (int, int, int, time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.totalSubmitted, e.totalCorrect, e.totalIncorrect, e.discoveryTime
}
