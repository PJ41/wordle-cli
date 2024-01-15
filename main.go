package main

import (
    "fmt"
    "wordle-cli/pkg/dictionary"
    "wordle-cli/pkg/term"
    "wordle-cli/pkg/game"
)

func main() {
    if err := dictionary.InitializationError; err != nil {
        fmt.Println("Error with dicitionary: ", err)
        return
    }

    if err := term.Configure(); err != nil {
        fmt.Println("Error configuring terminal: ", err)
    }
    defer term.Restore()

    game.Play()
}
