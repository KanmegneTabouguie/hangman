package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var words []string
var selectedWord string
var guessedWord []string
var lives int

func main() {
	lives = 10
	readWordsFromFile()
	selectRandomWord()
	initializeGuessedWord()

	for {
		displayGameStatus()

		if lives == 0 {
			fmt.Println("\nGame Over! The word was:", selectedWord)
			break
		}

		if !containsUnderscore(guessedWord) {
			fmt.Println("\nCongratulations! You guessed the word:", selectedWord)
			break
		}

		makeGuess()
	}
}

func readWordsFromFile() {
	file, err := os.Open("words.txt")
	if err != nil {
		panic("Error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
}

func selectRandomWord() {
	if len(words) == 0 {
		fmt.Println("No words available. Please check your words.txt file.")
		os.Exit(1)
	}

	selectedWord = words[rand.Intn(len(words))]
}

func initializeGuessedWord() {
	guessedWord = make([]string, len(selectedWord))
	for i := range guessedWord {
		guessedWord[i] = "_"
	}
}

func displayGameStatus() {
	fmt.Printf("\nWord: %s\n", strings.Join(guessedWord, " "))
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
	updateGuessedWord(guess)
}

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
