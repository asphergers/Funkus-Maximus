package discord

import (
    "github.com/bwmarrin/discordgo"
)


type BotCommand func(*discordgo.Session, *discordgo.MessageCreate)

type Song struct {
    url string
    title string
}

type Command struct {
    Name string
    Description string
    Help string
    Function BotCommand
}

