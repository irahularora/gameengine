package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game-engine/internal/models"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	n := 1000 // Scaled to 1000 users as requested
	fmt.Printf("Starting Mock User Engine with %d users...\n", n)

	var wg sync.WaitGroup
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// 1. Add a random delay (100–1000ms) to simulate network lag
			delay := time.Duration(rand.Intn(901)+100) * time.Millisecond
			time.Sleep(delay)

			// 2. Assign a correct answer flag (yes/no)
			randomAnswer := models.AnswerType(models.AnswerYes)
			if rand.Intn(2) == 0 {
				randomAnswer = models.AnswerNo
			}
			resp := models.UserResponse{
				UserID: fmt.Sprintf("user_%d", userID),
				Answer: randomAnswer,
			}

			// 3. Send response concurrently to the API server
			jsonData, _ := json.Marshal(resp)
			url := "http://localhost:8080/submit"

			_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				// Suppress individual errors for clean console during massive simulation
				// but log if it's a connection issue.
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All 1000 users have finished sending responses.")
}
