package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Test = Command {
    Name: "test",
    Description: "command for testing argument parsing",
    Help: "!test [string of arguments]",
    Function: test,
}

func test(s *discordgo.Session, m *discordgo.MessageCreate) {
    arguments := strings.SplitN(m.Content, " ", 2)[1]

    s.ChannelMessageSend(m.ChannelID, arguments);
}

