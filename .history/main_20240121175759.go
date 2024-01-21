package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var words map[string][]string
var selectedWord string
var guessedWord []string
var guessedLetters []string
var lives int
var score int

func main() {
	playHangman()
}

func playHangman() {
	lives = 10
	score = 0
	readWordsFromFile()

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
			fmt.Printf("\nCongratulations! You guessed the word: %s\n", selectedWord)
			updateScore(time.Since(startTime).Seconds())
			break
		}

		makeGuess()
	}

	fmt.Printf("Your score: %d\n", score)

	if playAgain() {
		playHangman()
	} else {
		fmt.Println("Thanks for playing Hangman!")
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
	guessedWord = make([]string, len(selectedWord))
	for i := range guessedWord {
		guessedWord[i] = "_"
	}

	guessedLetters = []string{}
}

func displayGameStatus() {
	fmt.Printf("\nWord: %s\n", strings.Join(guessedWord, " "))
	fmt.Printf("Guessed Letters: %s\n", strings.Join(guessedLetters, " "))
	fmt.Printf("Lives remaining: %d\n", lives)
}

func makeGuess() {
	var guess string
	fmt.Print("Enter a letter: ")
	fmt.Scanln(&guess)

	if len(guess) != 1 || !isLetter(guess) {
		fmt.Println("Please enter a valid single letter.")
		return
	}

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
		score += 10 // Increase score for correct guesses
	} else {
		fmt.Printf("Incorrect guess: %s\n", guess)
		lives--
		score -= 5 // Decrease score for incorrect guesses
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

func updateScore(seconds float64) {
	// Additional scoring based on time taken
	// You can customize this scoring logic according to your requirements
	score += int(100 / seconds)
}

func playAgain() bool {
	var response string
	fmt.Print("Do you want to play again? (yes/no): ")
	fmt.Scanln(&response)
	return strings.ToLower(response) == "yes"
}
