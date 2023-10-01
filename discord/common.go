package discord

import (
	"io"

	"github.com/bwmarrin/discordgo"
)


type BotCommand func(*discordgo.Session, *discordgo.MessageCreate)

type Song struct {
    url string
    title string
    buff *io.PipeReader
}

type Command struct {
    Name string
    Description string
    Help string
    Function BotCommand
}

