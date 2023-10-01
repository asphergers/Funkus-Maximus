package main

import (
    "fmt" 
    "main/discord"
    "main/parser"
)


func main() {
    fmt.Println("hello world")
    parser.InitParser()
    discord.Start()
}

