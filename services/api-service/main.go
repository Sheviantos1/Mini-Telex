// services/api-service/main.go
// This simulates telex_be - a simple REST API in Golang

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

// Simple in-memory storage
var users = make(map[string]User)
var messages = make(map[string]Message)

func main() {
	log.Println("Starting API Service on :8000")

	// Routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/users", usersHandler)
	http.HandleFunc("/api/messages", messagesHandler)
	http.HandleFunc("/api/chaos", chaosHandler) // For testing errors

	// Start server
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Health check requested")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "api-service",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// Users endpoint
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("[INFO] Fetching users list")
		w.Header().Set("Content-Type", "application/json")
		usersList := make([]User, 0, len(users))
		for _, user := range users {
			usersList = append(usersList, user)
		}
		json.NewEncoder(w).Encode(usersList)

	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("[ERROR] Failed to decode user: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user.ID = fmt.Sprintf("user_%d", rand.Intn(10000))
		user.CreatedAt = time.Now()
		users[user.ID] = user

		log.Printf("[INFO] User created successfully: %s", user.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Messages endpoint
func messagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("[INFO] Fetching messages list")
		w.Header().Set("Content-Type", "application/json")
		messagesList := make([]Message, 0, len(messages))
		for _, msg := range messages {
			messagesList = append(messagesList, msg)
		}
		json.NewEncoder(w).Encode(messagesList)

	case "POST":
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Printf("[ERROR] Failed to decode message: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		msg.ID = fmt.Sprintf("msg_%d", rand.Intn(10000))
		msg.Status = "pending"
		messages[msg.ID] = msg

		log.Printf("[INFO] Message created: %s for user %s", msg.ID, msg.UserID)

		// Simulate random errors (10% chance)
		if rand.Float32() < 0.1 {
			log.Printf("[ERROR] Failed to process message %s - random error", msg.ID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Chaos endpoint - deliberately causes errors for testing
func chaosHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")

	switch action {
	case "slow":
		log.Println("[WARN] Chaos mode: Slow response triggered")
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Slow response completed"))

	case "error":
		log.Println("[ERROR] Chaos mode: Error triggered")
		http.Error(w, "Chaos error", http.StatusInternalServerError)

	case "crash":
		log.Println("[FATAL] Chaos mode: Crash triggered")
		panic("Chaos crash!")

	default:
		log.Println("[INFO] Chaos mode: Status check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Chaos mode available. Use ?action=slow|error|crash"))
	}
}// Testing PR automation

// Test endpoint for PR automation
func testHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("[INFO] Test endpoint called")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"test":"success"}`))
}
// CalculateSum adds two numbers
func CalculateSum(a, b int) int {
    return a + b
}
// test change
