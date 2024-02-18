package main

import (
    "bufio"
    "os"
    "fmt"
    "wordle-cli/pkg/dictionary"
    "wordle-cli/pkg/term"
    "wordle-cli/pkg/game"
    "wordle-cli/pkg/data"
)

func main() {
    if handleCommandLineArguments() {
        return
    }

    if err := dictionary.Init(); err != nil {
        fmt.Println("Error with dicitionary: ", err)
        return
    }

    if err := data.Init(); err != nil {
        fmt.Println("Error with loading user data: ", err)
        return
    }
    defer data.Close()

    if err := term.Configure(); err != nil {
        fmt.Println("Error configuring terminal: ", err)
        return
    }
    defer term.Restore()

    game.Play()
}

func handleCommandLineArguments() bool {
    if len(os.Args) == 1 {
        return false
    }

    if len(os.Args) > 2 {
        printHelp()
        return true
    }

    switch os.Args[1] {
    case "--clean":
        handleClean()
    default:
        printHelp()
    }

    return true
}

func printHelp() {
    fmt.Println(
`
HELP PAGE: wordle cli
To play, run without any arguments

Command line options:
    -h,
    --help,
    default -> prints this screen

    --clean -> deletes data and directory that was created by the program
               this deletion is based on the data directory
               in general, you should clear before changing the env variable


User data is stored in the standardized data directory for mac, linux, or windows
You can customize this behavior by setting the WORDLE_CLI_DATA_DIR env variable

Default data directories:
    Mac -> $HOME/Library/Application Support/
    Windows -> $LOCALAPPDATA/
    Linux -> $HOME/.config/
`)
}

func handleClean() {
    fmt.Print("Are you sure you want to delete user data [y/n] (default no): ")

    reader := bufio.NewReader(os.Stdin)
    char, _, err := reader.ReadRune()
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }

    switch char {
    case 'y', 'Y':
        if err := data.CleanUp(); err != nil {
            fmt.Println("Error cleaning up: ", err)
        } else {
            fmt.Println("Deleted user data.")
        }
    default:
        fmt.Println("Did not delete user data.")
    }
}
