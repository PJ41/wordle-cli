package constants

type GameStatusEnum int

const (
    Playing GameStatusEnum = iota
    Won 
    Lost
    Exiting
)

const WordLength int = 5
const MaxAttempts int = 6
