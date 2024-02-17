package data 

import (
    "runtime"
    "os"
    "path/filepath"
    "errors"
    "wordle-cli/pkg/constants"
    "wordle-cli/pkg/dictionary"
    "bufio"
    "strconv"
    "strings"
)

const (
    Wins = 0
    Losses = 1
    Streak = 2
    MaxStreak = 3
    GuessDistribution = 4
    StatsCount = constants.MaxAttempts + GuessDistribution
    
    wordIdx = 0
    gameStatus = 1
    metaCount = 2

    gameStateCount = constants.MaxAttempts * constants.WordLength

    appName = "wordle_cli"
)

type data struct {
    Stats [StatsCount]int
    GameMeta [metaCount]int
    GameState [constants.MaxAttempts][constants.WordLength]rune
}

var file *os.File
var userData data;

func Init() error {
    dataDir, err := getDataDirectory()
    if err != nil {
        return err
    }

    if err := os.Mkdir(dataDir, 0755); err != nil && !os.IsExist(err) {
        return err
    }

    fileName := getFileName(dataDir)
    
    localFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return err
    }

    file = localFile

    if err := parseCsvFile(); err != nil {
        return err
    }

    return nil
}

func CleanUp() error {
    dataDir, err := getDataDirectory()
    if err != nil {
        return err
    }

    if _, err := os.Stat(dataDir); os.IsNotExist(err) {
        return nil
	} else {
        fileName := getFileName(dataDir)
        if err := os.Remove(fileName); err != nil {
            return err
        }

        if err := os.Remove(dataDir); err != nil {
            return err
        }
	}

    return nil
}

func Close() {
    file.Close()
}

func GetData(gameState *[constants.MaxAttempts][constants.WordLength]rune, status *constants.GameStatusEnum, stats *[StatsCount]int) int {
    *gameState = userData.GameState
    *status = constants.GameStatusEnum(userData.GameMeta[gameStatus])
    *stats = userData.Stats

    for r, row := range gameState {
        if equalZero(&row) {
            return r
        }
    }

    return constants.MaxAttempts
}

func SyncData(gameState *[constants.MaxAttempts][constants.WordLength]rune, wordIndex int, status constants.GameStatusEnum, stats *[StatsCount]int) {
    userData.GameState = *gameState
    userData.GameMeta[wordIdx] = wordIndex
    userData.GameMeta[gameStatus] = int(status)
    userData.Stats = *stats
    saveData()
}

func saveData() {
    file.Truncate(0)
    file.Seek(0, 0)

    writer := bufio.NewWriter(file)

    var statData [StatsCount]string
    for i := range statData {
        statData[i] = strconv.Itoa(userData.Stats[i])
    }
    writer.WriteString(strings.Join(statData[:], ",") + "\n")

    var gameMetaData [metaCount]string
    for i := range gameMetaData {
        gameMetaData[i] = strconv.Itoa(userData.GameMeta[i])
    }
    writer.WriteString(strings.Join(gameMetaData[:], ",") + "\n")

    var gameStateData [gameStateCount]string
    for i := range gameStateData {
        ni, nj := i / constants.WordLength, i % constants.WordLength
        gameStateData[i] = strconv.Itoa(int(userData.GameState[ni][nj]))
    }
    writer.WriteString(strings.Join(gameStateData[:], ",") + "\n")

    writer.Flush()
}

func equalZero(row *[constants.WordLength]rune) bool {
    for _, letter := range row {
        if letter != 0 {
            return false
        }
    }

    return true
}

func parseCsvFile() error {
    setDefaultUserData()

    scanner := bufio.NewScanner(file)

    if !scanner.Scan() {
        return nil
    }

    line := scanner.Text()
    row := strings.Split(line, ",")
    if len(row) != StatsCount {
        return corruptData()
    }

    for i, str := range row {
        val, err := strconv.Atoi(str)
        if err != nil {
            setDefaultUserData()
            return err
        }
        userData.Stats[i] = val
    }

    if !scanner.Scan() {
        return corruptData()
    }

    line = scanner.Text()
    row = strings.Split(line, ",")
    if len(row) != metaCount {
        return corruptData()
    }

    if val, err := strconv.Atoi(row[wordIdx]); err != nil {
        setDefaultUserData()
        return err
    } else if val != dictionary.GetWordOfDayIndex() {
        return nil
    }

    if val, err := strconv.Atoi(row[gameStatus]); err != nil {
        setDefaultUserData()
        return err
    } else {
        userData.GameMeta[gameStatus] = val
    }

    if !scanner.Scan() {
        return corruptData()
    }

    line = scanner.Text()
    row = strings.Split(line, ",")
    if len(row) != gameStateCount {
        return corruptData()
    }

    for i, char := range row {
        val, err := strconv.Atoi(char)
        if err != nil {
            setDefaultUserData()
            return err
        }
        userData.GameState[i / constants.WordLength][i % constants.WordLength] = rune(val)
    }

    return nil
}

func corruptData() error {
    return errors.New("Corrupt data")
}

func setDefaultUserData() {
    userData = data {
        Stats: [StatsCount]int{},
        GameMeta: [metaCount]int{},
        GameState: [constants.MaxAttempts][constants.WordLength]rune{},
    }
}

func getDataDirectory() (string, error) {
    dataDir := os.Getenv("WORDLE_CLI_DATA_DIR")
    if dataDir == "" {
        switch platform := runtime.GOOS; platform {
        case "darwin":
            home := os.Getenv("HOME")
            if home == "" {
                return "", errors.New("Missing HOME environment variable")
            }
            dataDir = filepath.Join(home, "Library", "Application Support")
        case "windows":
            dataDir = os.Getenv("APPDATA")
            if dataDir == "" {
                return "", errors.New("Missing APPDATA environment variable")
            }
        case "linux":
            home := os.Getenv("HOME")
            if home == "" {
                return "", errors.New("Missing HOME environment variable")
            }
            dataDir = filepath.Join(home, ".config")
        default:
            return "", errors.New("Unsupported OS, must set WORDLE_CLI_DATA_DIR env variable")
        }
    }

    return filepath.Join(dataDir, appName), nil
}

func getFileName(dataDir string) string {
    return filepath.Join(dataDir, "wordle_cli_user_data.csv")
}
