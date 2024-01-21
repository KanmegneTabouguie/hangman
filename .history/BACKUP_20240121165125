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
var guessedLetters []string
var lives int

func main() {
	playHangman()
}

func playHangman() {
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

	if playAgain() {
		playHangman()
	} else {
		fmt.Println("Thanks for playing Hangman!")
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
	} else {
		fmt.Printf("Incorrect guess: %s\n", guess)
		lives--
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

func playAgain() bool {
	var response string
	fmt.Print("Do you want to play again? (yes/no): ")
	fmt.Scanln(&response)
	return strings.ToLower(response) == "yes"
}
