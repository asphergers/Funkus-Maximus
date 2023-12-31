package discord

import (
	"io"

	"github.com/bwmarrin/discordgo"
)


type BotCommand func(*discordgo.Session, *discordgo.MessageCreate)

type Song struct {
    url string
    title string
    length string
    buff *io.PipeReader
}

type Command struct {
    Name string
    Aliases []string
    Description string
    Help string
    Function BotCommand
}

