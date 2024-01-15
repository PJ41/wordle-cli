package game

import (
    "os"
    "wordle-cli/pkg/dictionary"
    "wordle-cli/pkg/term"
)

type gameStatusEnum int

const (
    PLAYING gameStatusEnum = iota
    WON
    LOST
    EXITING
)

const wordLen int = 5
const maxAttempts int = 6

var magicWord string
var gameState [maxAttempts][wordLen]rune
var gameStatus gameStatusEnum

func Play() {
    magicWord = dictionary.GetWordOfDay()

    term.Println("Welcome to Wordle")
    term.Println("Options = [ 1 : Quit ]")
    term.Println("")

    row := 0
    col := 0

    renderState(row)

    for {
        row, col = nextPress(row, col)

        if gameStatus == EXITING {
            return
        }

        reRenderState(row)
    }
}

func nextPress(row int, col int) (int, int) {
    var key [1]byte
	_, err := os.Stdin.Read(key[:])
	if err != nil {
        term.Println("Key press error: %s", err)
        gameStatus = EXITING
		return row, col
	}

    kp := rune(key[0])

    if kp == '1' {
        gameStatus = EXITING
    } else if gameStatus == WON || gameStatus == LOST {
        return row, col
    } else if isLetter(kp) && col < wordLen {
        gameState[row][col] = kp
        col++
    } else if isDelete(kp) && col > 0 {
        col--
        gameState[row][col] = 0
    } else if isEnter(kp) {
        slice := gameState[row][:]
        word := string(slice)
        if (dictionary.IsValidGuess(word)) {
            if word == magicWord {
                gameStatus = WON
                term.Println("Congratulations, you've won in %d %s!", row + 1, attemptFormatted(row))
            } else if row == maxAttempts - 1 {
                gameStatus = LOST
                term.Println("Sorry, you're out of attempts. The correct answer was %s.", magicWord)
            }

            col = 0
            row++
        }
    }

    return row, col
}

func isLetter(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isEnter(r rune) bool {
    return r == 13
}

func isDelete(r rune) bool {
    return r == 127 
}

func renderState(row int) {
    if gameStatus == WON || gameStatus == LOST {
        term.MoveCursorUp(1)
    }

    for i, word := range gameState {
		for j, letter := range word {
            if letter == 0 {
                term.Print(" [ ]")
            } else if i >= row {
                term.Print(" [%c]", letter)
            } else if letter == rune(magicWord[j]) {
                term.PrintColored(term.GREEN, " [%c]", letter)
            } else if magicWordContains(letter) {
                term.PrintColored(term.YELLOW, " [%c]", letter)
            } else {
                term.PrintColored(term.RED, " [%c]", letter)
            }
		}
		term.Println("")
	}

    if gameStatus == WON || gameStatus == LOST {
        term.MoveCursorDown(1)
    }

    term.Flush()
}

func magicWordContains(r rune) bool {
    for _, letter := range magicWord {
        if r == letter {
            return true
        }
    }
    return false
}

func reRenderState(row int) {
    term.MoveCursorUp(maxAttempts)
    renderState(row)
}

func attemptFormatted(row int) string {
    if row == 0 {
        return "attempt"
    }
    return "attempts"
}
