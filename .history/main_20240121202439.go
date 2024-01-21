package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" // Import the Gorilla Mux router
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

const leaderboardFile = "leaderboard.txt"

func main() {
	// Load words from file
	readWordsFromFile()

	// Start the HTTP server
	http.HandleFunc("/init", handleInit)
	http.HandleFunc("/guess", handleGuess)

	// Serve static files (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("frontend")))

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
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
	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		difficulty = "medium" // Default difficulty
	}

	gameState := initializeGame(difficulty)
	json.NewEncoder(w).Encode(gameState)
}

// Handler for processing a guess
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

	// Encode the updated game state and send it as the response
	gameState := GameState{
		Word:           strings.Join(guessedWord, " "),
		GuessedLetters: guessedLetters,
		Lives:          lives,
	}
	json.NewEncoder(w).Encode(gameState)
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
		initializeGame("medium") // Reset the game for the next round
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
