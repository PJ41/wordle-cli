package dictionary

import (
    "bufio"
    "os"
    "path/filepath"
    "runtime"
    "time"
    "strings"
)

const answersLineCount = 2315
const answersFileName string = "wordle-answers-shuffled.txt"
const nonAnswersFileName string = "wordle-non-answers-sorted.txt"

var answers map[string]bool
var keys []string
var nonAnswers map[string]bool

func Init() error {
    answers = make(map[string]bool)
    keys = make([]string, 0, answersLineCount)
    nonAnswers = make(map[string]bool)

    if err := loadWordsIntoMap(answers, answersFileName); err != nil {
        return err
    }

    if err := loadWordsIntoMap(nonAnswers, nonAnswersFileName); err != nil {
        return err
    }

    return nil
}

func GetWord(wordIndex int) string {
    return strings.ToUpper(keys[wordIndex])
}

func GetWordOfDayIndex() int {
    referenceTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
    currentTime := time.Now()

    durationSince := currentTime.Sub(referenceTime)
    daysSince := int(durationSince.Hours() / 24)

    return daysSince % len(keys) 
}

func IsValidGuess(guess string) bool {
    guess = strings.ToLower(guess)
    return answers[guess] || nonAnswers[guess]
}

func loadWordsIntoMap(theMap map[string]bool, fileName string) error {
    _, currentFile, _, _ := runtime.Caller(0)
    dir := filepath.Dir(currentFile)

    filePath := filepath.Join(dir, "words", fileName)

    file, err := os.Open(filePath); 
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        theMap[line] = true

        if fileName == answersFileName {
            keys = append(keys, line)
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}
