package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Skip = Command {
    Name: "skip",
    Aliases: []string{"skip", "s"},
    Description: "play audio from youtube videos",
    Help: "!play [url]",
    Function: skip,
}

func skip(s *discordgo.Session, m *discordgo.MessageCreate) {
    guild, ok := GuildMap[m.GuildID]
    if !ok {
        returnMessage := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }
    
    if len(guild.Queue) == 0 && guild.CurrentStream == nil {
        returnMessage := fmt.Sprintf("queue is empty")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    guild.CurrentStream.Kill()
}
