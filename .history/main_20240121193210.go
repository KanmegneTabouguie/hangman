package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

var words map[string][]string
var selectedWord string
var guessedWord []string
var guessedLetters []string
var lives int
var score int
var currentPlayer int
var players []string
var playerScores map[string]int
var guessTimeout time.Duration
var db *sql.DB

const (
	leaderboardTable = "hangman_scores"
)

func main() {
	setupDatabase()
	defer db.Close()

	playHangman()
}

func setupDatabase() {
	connStr := "user=username dbname=yourdb sslmode=disable" // Replace with your actual database credentials
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createLeaderboardTable()
}

func createLeaderboardTable() {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			player VARCHAR(255) PRIMARY KEY,
			score INT
		)
	`, leaderboardTable))
	if err != nil {
		log.Fatal(err)
	}
}

func playHangman() {
	lives = 10
	score = 0
	currentPlayer = 0
	playerScores = make(map[string]int)
	readWordsFromFile()

	// Get the number of players
	numPlayers := selectNumPlayers()

	// Get player names
	getPlayerNames(numPlayers)

	// Set the guess timeout
	setGuessTimeout()

	difficulty := selectDifficulty()
	selectRandomWord(difficulty)
	initializeGuessedWord()

	startTime := time.Now()

	for {
		displayGameStatus()

		if lives == 0 {
			fmt.Println("\nGame Over! The word was:", selectedWord)
			break
		}

		if !containsUnderscore(guessedWord) {
			winner := determineWinner()
			fmt.Printf("\nCongratulations! %s guessed the word: %s\n", winner, selectedWord)
			updateScore(time.Since(startTime).Seconds(), winner)
			break
		}

		makeGuessWithTimeout()

		// Switch to the next player in multiplayer mode
		if numPlayers > 1 {
			currentPlayer = (currentPlayer + 1) % numPlayers
		}
	}

	displayScores()

	updateLeaderboard()

	if playAgain() {
		playHangman()
	} else {
		fmt.Println("Thanks for playing Hangman!")
	}
}

func setGuessTimeout() {
	fmt.Print("Enter the guess time limit in seconds (e.g., 30): ")
	fmt.Scanln(&guessTimeout)
}

func selectNumPlayers() int {
	var numPlayers int
	for {
		fmt.Print("Enter the number of players (1, 2, or 3): ")
		fmt.Scanln(&numPlayers)
		if numPlayers >= 1 && numPlayers <= 3 {
			break
		}
		fmt.Println("Invalid number of players. Please enter 1, 2, or 3.")
	}
	return numPlayers
}

func getPlayerNames(numPlayers int) {
	players = make([]string, numPlayers)
	for i := 0; i < numPlayers; i++ {
		fmt.Printf("Enter the name of player %d: ", i+1)
		fmt.Scanln(&players[i])
		playerScores[players[i]] = 0
	}
}

func readWordsFromFile() {
	words = make(map[string][]string)
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
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
}

func selectDifficulty() string {
	var difficulty string
	fmt.Print("Select difficulty level (easy/medium/hard): ")
	fmt.Scanln(&difficulty)
	return difficulty
}

func selectRandomWord(difficulty string) {
	wordList, ok := words[difficulty]
	if !ok || len(wordList) == 0 {
		fmt.Println("No words available for the selected difficulty. Please check your words.txt file.")
		os.Exit(1)
	}

	selectedWord = wordList[rand.Intn(len(wordList))]
}

func initializeGuessedWord() {
	guessedWord = make([]string, len(selectedWord))
	for i := range guessedWord {
		guessedWord[i] = "_"
	}

	guessedLetters = []string{}
}

func displayGameStatus() {
	fmt.Printf("\nCurrent Player: %s\n", players[currentPlayer])
	fmt.Printf("Word: %s\n", strings.Join(guessedWord, " "))
	fmt.Printf("Guessed Letters: %s\n", strings.Join(guessedLetters, " "))
	fmt.Printf("Lives remaining: %d\n", lives)
}

func makeGuessWithTimeout() {
	ch := make(chan string, 1)
	go func() {
		makeGuess()
		ch <- "done"
	}()
	select {
	case <-ch:
		// Guess completed within the time limit
	case <-time.After(time.Second * guessTimeout):
		// Time limit exceeded
		fmt.Printf("Time limit exceeded! %s's turn is considered an incorrect guess.\n", players[currentPlayer])
		lives--
		playerScores[players[currentPlayer]] -= 5 // Penalize for exceeding time limit
	}
}

func makeGuess() {
	var guess string
	fmt.Printf("%s, enter a letter or type 'hint' for a hint: ", players[currentPlayer])
	fmt.Scanln(&guess)

	if guess == "hint" {
		revealRandomLetter()
	} else if len(guess) != 1 || !isLetter(guess) {
		fmt.Println("Please enter a valid single letter.")
		return
	} else {
		processGuess(guess)
	}
}

func revealRandomLetter() {
	// Get a list of unrevealed indices
	unrevealedIndices := []int{}
	for i, char := range guessedWord {
		if char == "_" {
			unrevealedIndices = append(unrevealedIndices, i)
		}
	}

	if len(unrevealedIndices) == 0 {
		fmt.Println("No more letters to reveal.")
		return
	}

	// Randomly select an index and reveal the corresponding letter
	randomIndex := unrevealedIndices[rand.Intn(len(unrevealedIndices))]
	guessedWord[randomIndex] = string(selectedWord[randomIndex])

	fmt.Printf("Hint: The letter at position %d is %s\n", randomIndex+1, guessedWord[randomIndex])
}

func processGuess(guess string) {
	guess = strings.ToLower(guess)

	if contains(guessedLetters, guess) {
		fmt.Printf("You already guessed the letter %s. Try another one.\n", guess)
		return
	}

	guessedLetters = append(guessedLetters, guess)
	updateGuessedWord(guess)
}

func updateGuessedWord(guess string) {
	if strings.Contains(selectedWord, guess) {
		for i, char := range selectedWord {
			if string(char) == guess {
				guessedWord[i] = guess
			}
		}
		playerScores[players[currentPlayer]] += 10 // Increase score for correct guesses
	} else {
		fmt.Printf("Incorrect guess: %s\n", guess)
		lives--
		playerScores[players[currentPlayer]] -= 5 // Decrease score for incorrect guesses
	}
}

func contains(array []string, value string) bool {
	for _, elem := range array {
		if elem == value {
			return true
		}
	}
	return false
}

func containsUnderscore(word []string) bool {
	for _, char := range word {
		if char == "_" {
			return true
		}
	}
	return false
}

func isLetter(s string) bool {
	return len(s) == 1 && ('a' <= s[0] && s[0] <= 'z' || 'A' <= s[0] && s[0] <= 'Z')
}

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

func displayScores() {
	fmt.Println("\nFinal Scores:")
	for _, player := range players {
		fmt.Printf("%s: %d\n", player, playerScores[player])
	}
}

func updateScore(seconds float64, winner string) {
	// Additional scoring based on time taken
	// You can customize this scoring logic according to your requirements
	playerScores[winner] += int(100 / seconds)
}

func updateLeaderboard() {
	// Read existing leaderboard
	leaderboard := readLeaderboard()

	// Add current game scores to the leaderboard
	for _, player := range players {
		entry := LeaderboardEntry{Player: player, Score: playerScores[player]}
		leaderboard = append(leaderboard, entry)
	}

	// Sort the leaderboard by score in descending order
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	// Keep only the top 10 scores
	if len(leaderboard) > 10 {
		leaderboard = leaderboard[:10]
	}

	// Write the updated leaderboard to the database
	writeLeaderboardToDB(leaderboard)

	// Display the updated leaderboard
	fmt.Println("\nLeaderboard:")
	for i, entry := range leaderboard {
		fmt.Printf("%d. %s: %d\n", i+1, entry.Player, entry.Score)
	}
}

func readLeaderboard() []LeaderboardEntry {
	var leaderboard []LeaderboardEntry
	rows, err := db.Query(fmt.Sprintf("SELECT player, score FROM %s ORDER BY score DESC LIMIT 10", leaderboardTable))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var player string
		var score int
		if err := rows.Scan(&player, &score); err != nil {
			log.Fatal(err)
		}
		entry := LeaderboardEntry{Player: player, Score: score}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard
}

func writeLeaderboardToDB(leaderboard []LeaderboardEntry) {
	// Clear existing entries in the leaderboard table
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", leaderboardTable))
	if err != nil {
		log.Fatal(err)
	}

	// Insert the new entries into the leaderboard table
	for _, entry := range leaderboard {
		_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (player, score) VALUES ($1, $2)", leaderboardTable), entry.Player, entry.Score)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// LeaderboardEntry represents a player's entry in the leaderboard.
type LeaderboardEntry struct {
	Player string
	Score  int
}

func playAgain() bool {
	var response string
	fmt.Print("Do you want to play again? (yes/no): ")
	fmt.Scanln(&response)
	return strings.ToLower(response) == "yes"
}
