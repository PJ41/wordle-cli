package game

import (
    "os"
    "wordle-cli/pkg/dictionary"
    "wordle-cli/pkg/term"
    "wordle-cli/pkg/constants"
    "wordle-cli/pkg/data"
    "time"
    "unicode"
)

const enterKey rune = 13

var wordIndex int
var magicWord string
var gameState [constants.MaxAttempts][constants.WordLength]rune
var gameStatus constants.GameStatusEnum
var stats [data.StatsCount]int
var currentScreen screen
var offset = 0;

type screen interface {
    screenNumber() int
    render(row int)
    rePositionCursor()
    clear()
}

type playScreen struct{}

func (ps playScreen) screenNumber() int {
    return 1
}

func (ps playScreen) render(row int) {
    for i, word := range gameState {
        yellows := make(map[rune]int)
        for j, letter := range magicWord {
            if letter != word[j] {
                yellows[letter]++
            }
        }

        for j, letter := range word {
            if letter == 0 {
                term.Print(" [ ]")
            } else if i >= row {
                term.Print(" [%c]", letter)
            } else if letter == rune(magicWord[j]) {
                term.PrintModified(term.GREEN, " [%c]", letter)
            } else if yellows[letter] > 0 {
                yellows[letter]--
                term.PrintModified(term.YELLOW, " [%c]", letter)
            } else {
                term.PrintModified(term.RED, " [%c]", letter)
            }
        }
        term.Println("")
    }

    if gameStatus == constants.Won {
        printVictory(row)
    } else if gameStatus == constants.Lost {
        printDefeat()
    }

    term.Flush()
}

func (ps playScreen) rePositionCursor() {
    toMove := constants.MaxAttempts
    if gameStatus == constants.Won || gameStatus == constants.Lost {
        toMove++
    }
    term.MoveCursorUp(toMove)
}

func (ps playScreen) clear() {
    toDelete := constants.MaxAttempts
    if gameStatus == constants.Won || gameStatus == constants.Lost {
        toDelete++
    }
    term.ClearRows(toDelete)
}

type statsScreen struct{}

func (ss statsScreen) screenNumber() int {
    return 2
}

func (ss statsScreen) render(row int) {
    total := stats[data.Wins] + stats[data.Losses]

    term.Println("Played: %d", total)
    if total != 0 {
        percentage := int(float64(stats[data.Wins]) / float64(total) * 100)
        term.Println("Win %%: %d", percentage)
    } else {
        term.Println("Win %%: N/A")
    }
    term.Println("Current Streak: %d", stats[data.Streak])
    term.Println("Max Streak: %d", stats[data.MaxStreak])
    term.Println("Guess Distribution")

    for i := 0; i < constants.MaxAttempts; i++ {
        term.Println(" %d: %d", i + 1, stats[data.GuessDistribution + i])
    }

    term.Flush()
}

func (ss statsScreen) rePositionCursor() {
    term.MoveCursorUp(data.StatsCount + 1)
}

func (ss statsScreen) clear() {
    term.ClearRows(data.StatsCount + 1)
}

func Play() {
    currentScreen = playScreen{}

    wordIndex = dictionary.GetWordOfDayIndex()
    magicWord = dictionary.GetWord(wordIndex)

    currentTime := time.Now()
    term.Println("Welcome to Wordle")
    term.Println("Today is: %s", currentTime.Format("January 2, 2006"))
    printMenu()

    row := data.GetData(&gameState, &gameStatus, &stats); 
    col := 0

    currentScreen.render(row)

    for {
        row, col = nextPress(row, col)

        if gameStatus == constants.Exiting {
            return
        }

        currentScreen.rePositionCursor()
        currentScreen.render(row)
    }
}

func nextPress(row int, col int) (int, int) {
    var key [1]byte
    _, err := os.Stdin.Read(key[:])
    if err != nil {
        term.Println("Key press error: %s", err)
        gameStatus = constants.Exiting
        return row, col
    }

    kp := rune(key[0])
    return processKey(kp, row, col)
}

func processKey(kp rune, row int, col int) (int, int) {
    if kp == '3' {
        gameStatus = constants.Exiting
    } else if kp == '1' {
        swapScreens(playScreen{}, row)
    } else if kp == '2' {
        swapScreens(statsScreen{}, row)
    } else if _, ok := currentScreen.(statsScreen); ok {
        return row, col
    } else if gameStatus == constants.Won || gameStatus == constants.Lost {
        return row, col
    } else if isLetter(kp) && col < constants.WordLength {
        kp = unicode.ToUpper(kp)
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
                stats[data.Wins]++
                stats[data.Streak]++
                stats[data.MaxStreak] = max(stats[data.MaxStreak], stats[data.Streak])
                stats[data.GuessDistribution + row]++
                term.Println("")
                gameStatus = constants.Won
            } else if row == constants.MaxAttempts - 1 {
                stats[data.Losses]++
                stats[data.Streak] = 0
                term.Println("")
                gameStatus = constants.Lost
            }

            data.SyncData(&gameState, wordIndex, gameStatus, &stats)

            col = 0
            row++
        }
    }

    return row, col
}

func printVictory(row int) {
    term.Println("Congratulations, you've won in %d %s!", row, attemptFormatted(row))
}

func printDefeat() {
    term.Println("Sorry, you're out of attempts. The correct answer was %s.", magicWord)
}

func isLetter(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isEnter(r rune) bool {
    return r == enterKey
}

func isDelete(r rune) bool {
    return r == 127 
}

func printMenu() {
    term.Print("Menu { 1: ")
    if _, ok := currentScreen.(playScreen); ok {

        term.PrintModified(term.UNDERLINED, "Play")
    } else {
        term.Print("Play")
    }
    term.Print(", 2: ")

    if _, ok := currentScreen.(statsScreen); ok {
        term.PrintModified(term.UNDERLINED, "Stats")
    } else {
        term.Print("Stats")
    }

    term.Println(", 3: Quit }")
    term.Println("")
}

func swapScreens(newScreen screen, row int) {
    if (currentScreen.screenNumber() == newScreen.screenNumber()) {
        return
    }

    currentScreen.clear()
    currentScreen = newScreen

    term.MoveCursorUp(2)
    printMenu()

    currentScreen.render(row)
}

func attemptFormatted(row int) string {
    if row == 0 {
        return "attempt"
    }
    return "attempts"
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
