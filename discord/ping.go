package discord

import (
    "github.com/bwmarrin/discordgo"
)

var Ping = Command {
    Name: "ping",
    Description: "standard ping command",
    Help: "!ping",
    Function: ping,
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
    s.ChannelMessageSend(m.ChannelID, "pong")
}
