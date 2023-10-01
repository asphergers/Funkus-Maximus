package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Queue = Command {
    Name: "q",
    Description: "show the current queue",
    Help: "!play [url]",
    Function: queue,
}

func queue(s *discordgo.Session, m *discordgo.MessageCreate) {
    guild, ok := GuildMap[m.GuildID]
    if !ok {
        returnMessage := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    if len(guild.Queue) == 0 {
        returnMessage := fmt.Sprintf("queue is empty")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    var message string
    for _, song := range guild.Queue {
        entry := song.title + "\n" 
        message += entry
    }

    s.ChannelMessageSend(m.ChannelID, message)
}
