package main

import (
    "fmt"
    "wordle-cli/pkg/dictionary"
    "wordle-cli/pkg/term"
    "wordle-cli/pkg/game"
    "wordle-cli/pkg/data"
)

func main() {
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
