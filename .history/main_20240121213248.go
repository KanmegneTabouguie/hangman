package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	//"github.com/gorilla/mux" // Import the Gorilla Mux router
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync" // Import the sync package
	"time"
)

var words map[string][]string
var selectedWord string
var guessedWord []string
var guessedLetters []string
var lives int
var currentPlayer int
var players []string
var playerScores map[string]int
var guessTimeout time.Duration
var currentPlayerIndex int // Added to keep track of the current player's index

const leaderboardFile = "leaderboard.txt"

func main() {
	// Load words from file
	readWordsFromFile()

	// Initialize playerScores map
	playerScores = make(map[string]int)
	// Create a WaitGroup to wait for the server goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the HTTP server in a goroutine
	go func() {
		defer wg.Done()
		http.HandleFunc("/init", handleInit)
		http.HandleFunc("/guess", handleGuess)

		// Serve static files (HTML, CSS, JS)
		http.Handle("/", http.FileServer(http.Dir("public")))

		fmt.Println("Server listening on :8083")
		if err := http.ListenAndServe(":8083", nil); err != nil {
			fmt.Println("Error starting the server:", err)
		}
	}()

	// Wait for the server goroutine to finish
	wg.Wait()
}

// Struct for API requests and responses
type GuessData struct {
	Guess string `json:"guess"`
}

type GameState struct {
	Word           string   `json:"word"`
	GuessedLetters []string `json:"guessedLetters"`
	Lives          int      `json:"lives"`
}

// Handler for initializing a new game
func handleInit(w http.ResponseWriter, r *http.Request) {
	// Extract player name from query parameters
	playerName := r.URL.Query().Get("playerName")
	if playerName == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}

	// Add the player to the list
	players = append(players, playerName)
	playerScores[playerName] = 0

	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		difficulty = "medium" // Default difficulty
	}

	gameState := initializeGame(difficulty)
	json.NewEncoder(w).Encode(gameState)
	currentPlayerIndex = 0 // Initialize the current player index
	// Set a default guess timeout (e.g., 30 seconds)
	guessTimeout = 30 * time.Second

}

// Handler for processing a guess
func handleGuess(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body
	var guessData GuessData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &guessData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if the player exists
	if currentPlayer >= len(players) {
		http.Error(w, "Invalid player", http.StatusBadRequest)
		return
	}

	// Special case: Check for a hint request
	if guessData.Guess == "hint" {
		hint := provideHint()
		fmt.Printf("%s\n", hint)
		// Encode the updated game state and send it as the response
		gameState := GameState{
			Word:           strings.Join(guessedWord, " "),
			GuessedLetters: guessedLetters,
			Lives:          lives,
		}
		json.NewEncoder(w).Encode(gameState)
		return
	}

	// Update the game state based on the guess
	processGuess(guessData.Guess)

	// Switch to the next player's turn
	currentPlayerIndex = (currentPlayerIndex + 1) % len(players)
	currentPlayer = currentPlayerIndex

	// Set a timeout for the next player's turn
	timeoutTimer := time.NewTimer(guessTimeout)
	go func() {
		<-timeoutTimer.C
		handleTimeout()
	}()

	// Encode the updated game state and send it as the response
	gameState := GameState{
		Word:           strings.Join(guessedWord, " "),
		GuessedLetters: guessedLetters,
		Lives:          lives,
	}
	json.NewEncoder(w).Encode(gameState)
}

// Function to handle player timeout
func handleTimeout() {
	// Implement timeout logic here
	// For example, skip the current player's turn or end the game

	// For now, simply print a message
	fmt.Println("Player timeout!")
}

// Function to initialize a new game
func initializeGame(difficulty string) GameState {
	wordList, ok := words[difficulty]
	if !ok || len(wordList) == 0 {
		fmt.Println("No words available for the selected difficulty.")
		os.Exit(1)
	}

	selectedWord = wordList[rand.Intn(len(wordList))]
	guessedWord = make([]string, len(selectedWord))
	for i := range guessedWord {
		guessedWord[i] = "_"
	}
	guessedLetters = []string{}
	lives = 10

	// Reset player-specific data for a new game
	guessedWord = make([]string, len(selectedWord))
	for i := range guessedWord {
		guessedWord[i] = "_"
	}
	guessedLetters = []string{}
	lives = 10

	return GameState{
		Word:           strings.Join(guessedWord, " "),
		GuessedLetters: guessedLetters,
		Lives:          lives,
	}
}

func provideHint() string {
	// Implement your hint logic here
	// For example, you can reveal a random letter from the selectedWord
	// or provide a hint based on some other criteria

	// In this example, reveal a random letter
	hiddenIndices := make([]int, 0)
	for i, char := range guessedWord {
		if char == "_" {
			hiddenIndices = append(hiddenIndices, i)
		}
	}

	if len(hiddenIndices) == 0 {
		// No hidden letters to reveal
		return "No hint available."
	}

	randomIndex := hiddenIndices[rand.Intn(len(hiddenIndices))]
	guessedWord[randomIndex] = string(selectedWord[randomIndex])

	return fmt.Sprintf("Hint: Revealed letter at position %d.", randomIndex+1)
}

// Function to process a guess
func processGuess(guess string) {
	guess = strings.ToLower(guess)

	if contains(guessedLetters, guess) {
		fmt.Printf("You already guessed the letter %s. Try another one.\n", guess)
		return
	}

	guessedLetters = append(guessedLetters, guess)
	updateGuessedWord(guess)
}

// Function to update guessedWord based on the guess
func updateGuessedWord(guess string) {
	if strings.Contains(selectedWord, guess) {
		for i, char := range selectedWord {
			if string(char) == guess {
				guessedWord[i] = guess
			}
		}
	} else {
		fmt.Printf("Incorrect guess: %s\n", guess)
		lives--
	}

	if !contains(guessedWord, "_") {
		winner := determineWinner()
		fmt.Printf("\nCongratulations! %s guessed the word: %s\n", winner, selectedWord)
		updateScore(winner)
		saveToLeaderboard(winner, playerScores[winner])
		initializeGame("medium") // Reset the game for the next round
	}
}

// Function to save game results to the leaderboard file
func saveToLeaderboard(winner string, score int) {
	file, err := os.OpenFile(leaderboardFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening leaderboard file:", err)
		return
	}
	defer file.Close()

	result := fmt.Sprintf("%s: %d\n", winner, score)
	if _, err := file.WriteString(result); err != nil {
		fmt.Println("Error writing to leaderboard file:", err)
	}
}

// Function to determine the winner
func determineWinner() string {
	winner := players[0]
	highestScore := playerScores[players[0]]

	for _, player := range players[1:] {
		if playerScores[player] > highestScore {
			winner = player
			highestScore = playerScores[player]
		}
	}

	return winner
}

// Function to update scores based on time taken
func updateScore(winner string) {
	// TODO: Implement logic to update player scores
	// You might want to modify this based on your game rules
	playerScores[winner] += 10
}

// Function to read words from file
func readWordsFromFile() {
	words = make(map[string][]string)

	file, err := os.Open("words.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		parts := strings.Split(word, ":")
		if len(parts) == 2 {
			difficulty := parts[0]
			word := parts[1]
			words[difficulty] = append(words[difficulty], word)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
}

// Function to check if a value is in an array
func contains(array []string, value string) bool {
	for _, elem := range array {
		if elem == value {
			return true
		}
	}
	return false
}
