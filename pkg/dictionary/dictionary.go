package dictionary

import (
    "bufio"
    "os"
    "path/filepath"
    "runtime"
    "time"

)

var InitializationError error

const answersLineCount = 2315
const answersFileName string = "wordle-answers-shuffled.txt"
const nonAnswersFileName string = "wordle-non-answers-sorted.txt"

var answers map[string]bool
var keys []string
var nonAnswers map[string]bool

func init() {
    answers = make(map[string]bool)
    keys = make([]string, 0, answersLineCount)
    nonAnswers = make(map[string]bool)

    loadWordsIntoMap(answers, answersFileName)
    loadWordsIntoMap(nonAnswers, nonAnswersFileName)
}

func GetWordOfDay() string {
    durationSinceEpoch := time.Since(time.Unix(0, 0))
    daysSinceEpoch := int(durationSinceEpoch.Hours() / 24)

    index := daysSinceEpoch % len(keys) 
    return keys[index]
}

func IsValidGuess(guess string) bool {
    return answers[guess] || nonAnswers[guess]
}

func loadWordsIntoMap(theMap map[string]bool, fileName string) {
    if (InitializationError != nil) {
        return
    }

    _, currentFile, _, _ := runtime.Caller(0)
    dir := filepath.Dir(currentFile)

    filePath := filepath.Join(dir, "words", fileName)

    file, err := os.Open(filePath); 
    if err != nil {
        InitializationError = err
        return 
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
        InitializationError = err
    }
}
