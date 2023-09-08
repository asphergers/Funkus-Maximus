package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Skip = Command {
    Name: "skip",
    Description: "play audio from youtube videos",
    Help: "!play [url]",
    Function: skip,
}

func skip(s *discordgo.Session, m *discordgo.MessageCreate) {
    guild, ok := GuildMap[m.GuildID]
    if !ok {
        err := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    if len(guild.Queue) == 0 {
        err := fmt.Sprintf("queue is empty")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    guild.CurrentStream.Kill()
}
