package main

import (
	"bufio"
	"fmt"
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
var currentPlayer int
var players []string
var playerScores map[string]int
var guessTimeout time.Duration

const leaderboardFile = "leaderboard.txt"

func main() {
	playHangman()
}

func playHangman() {
	lives = 10
	currentPlayer = 0
	playerScores = make(map[string]int)
	readWordsFromFile()

	numPlayers := selectNumPlayers()
	getPlayerNames(numPlayers)
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
		panic("Error opening file")
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
	guessedWord = strings.Split(strings.Repeat("_", len(selectedWord)), "")
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
	playerScores[winner] += int(100 / seconds)
}

func updateLeaderboard() {
	leaderboard := readLeaderboard()

	for _, player := range players {
		entry := LeaderboardEntry{Player: player, Score: playerScores[player]}
		leaderboard = append(leaderboard, entry)
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	if len(leaderboard) > 10 {
		leaderboard = leaderboard[:10]
	}

	writeLeaderboard(leaderboard)

	fmt.Println("\nLeaderboard:")
	for i, entry := range leaderboard {
		fmt.Printf("%d. %s: %d\n", i+1, entry.Player, entry.Score)
	}
}

func readLeaderboard() []LeaderboardEntry {
	file, err := os.OpenFile(leaderboardFile, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening leaderboard file:", err)
		return nil
	}
	defer file.Close()

	var leaderboard []LeaderboardEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			player := parts[0]
			score := parseScore(parts[1])
			entry := LeaderboardEntry{Player: player, Score: score}
			leaderboard = append(leaderboard, entry)
		}
	}

	return leaderboard
}

func writeLeaderboard(leaderboard []LeaderboardEntry) {
	file, err := os.OpenFile(leaderboardFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening leaderboard file:", err)
		return
	}
	defer file.Close()

	for _, entry := range leaderboard {
		line := fmt.Sprintf("%s:%d\n", entry.Player, entry.Score)
		file.WriteString(line)
	}
}

func parseScore(scoreStr string) int {
	score, err := fmt.Sscanf(scoreStr, "%d")
	if err != nil {
		fmt.Println("Error parsing score:", err)
		return 0
	}
	return score
}

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
