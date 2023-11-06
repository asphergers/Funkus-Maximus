package discord

import (
    "fmt"

	"github.com/bwmarrin/discordgo"
)

var Clear = Command {
    Name: "clear",
    Aliases: []string{"clear"},
    Description: "clear the queue",
    Help: "!clear",
    Function: clear,
}

func clear(s *discordgo.Session, m *discordgo.MessageCreate) {
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

    guild.Queue = make([]*Song, 0)
}
