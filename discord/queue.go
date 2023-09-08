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
        err := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    if len(guild.Queue) == 0 {
        err := fmt.Sprintf("queue is empty")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    var message string

    for _, song := range guild.Queue {
        entry := song.title + "\n" 
        message += entry
    }

    s.ChannelMessageSend(m.ChannelID, message)
}
