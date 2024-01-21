package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var table [100]string
var letter string
var table2 []string
var x1 []byte
var y int
var lives int

func demonstration() {
	switch lives {
	case 9:
		fmt.Printf("\n \n \n \n \n \n=========\n")
	case 8:
		fmt.Printf("\n      |  \n      |  \n      |  \n      |  \n      |  \n=========\n")
	case 7:
		fmt.Printf("  +---+\n      |  \n      |  \n      |  \n      |  \n      |  \n=========\n")
	case 6:
		fmt.Printf("  +---+\n  |   |  \n      |  \n      |  \n      |  \n      |  \n=========\n")
	case 5:
		fmt.Printf("  +---+\n  |   |  \n  O   |  \n      |  \n      |  \n      |  \n=========\n")
	case 4:
		fmt.Printf("  +---+\n  |   |  \n  O   |  \n  |   |  \n      |  \n      |  \n=========\n")
	case 3:
		fmt.Printf("  +---+\n  |   |  \n  O   |  \n /|   |  \n      |  \n      |  \n=========\n")

	case 2:
		fmt.Print("  +---+\n  |   |  \n  O   |  \n /|", string(92), "  |  \n      |  \n      |  \n=========\n")
	case 1:
		fmt.Print("  +---+\n  |   |  \n  O   |  \n /|", string(92), "  |  \n /    |  \n      |  \n=========\n")
	case 0:
		fmt.Print("  +---+\n  |   |  \n  O   |  \n /|", string(92), "  |  \n / ", string(92), "  |  \n      |  \n=========\n")

	}
}

func lireMot() {
	content, err := os.Open("words.txt")

	if err != nil {
		panic("dommage erreur")
	}

	i := 0
	contentScanner := bufio.NewScanner(content)
	for contentScanner.Scan() {

		table[i] = contentScanner.Text()
		i++
	}

}

func tableFill() {
	table2 = make([]string, len(x1))
	for i := 0; i < len(table2); i++ {
		table2[i] = "_"
	}
}

func displayTable() {
	for i := 0; i < len(table2); i++ {
		fmt.Printf(table2[i])
		fmt.Printf("")
	}
}

func instruction() {
	fmt.Println("complet the following word below", lives, "lives")
	displayTable()
	fmt.Println("")
	fmt.Scanln(&letter)
	y = 0

	differentiation()
}

func splitWord(x int) {
	x1 = []byte(table[x])
	tableFill()
			lives--
		} else if letter == string(x1[h]) {
			for i := 0; i < len(x1); i++ {
				if letter == string(x1[i]) {
					table2[i] = letter
					continue
				}
			}
			break
		}
	}

	demonstration()

	fmt.Println("")
	if lives == 0 {
		fmt.Println(" ")
		fmt.Println("_|_|_|    _|_|_|_|  _|_|_|    _|_|_|    _|    _|")
		fmt.Println("_|    _|  _|        _|    _|  _|    _|  _|    _|")
		fmt.Println("_|_|_|    _|_|_|    _|_|_|    _|    _|  _|    _|")
		fmt.Println("_|        _|        _|    _|  _|    _|  _|    _|")
		fmt.Println("_|        _|_|_|_|  _|    _|  _|_|_|      _|_|  ")

		main()

	} else {
		for i := 0; i < len(table2); i++ {
			if table2[i] == "_" {
				instruction()
			} else if i == len(table2)-1 {
				displayTable()

				fmt.Println("")
				fmt.Println("_|          _| _|_|_|  _|      _|")
				fmt.Println("_|          _|   _|    _|_|    _|")
				fmt.Println("_|    _|    _|   _|    _|  _|  _|")
				fmt.Println("  _|  _|  _|     _|    _|    _|_|")
				fmt.Println("    _|  _|     _|_|_|  _|      _|")
				main()

			}
		}
	}
}

func main() {
	lives = 10
	y = rand.Intn(100)
	lireMot()
	splitWord(y)
}
