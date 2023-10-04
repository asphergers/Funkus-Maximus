package discord

import (
    "fmt"

	"github.com/bwmarrin/discordgo"
)

var Pause = Command {
    Name: "pause",
    Aliases: []string{"pause", "pa"},
    Description: "show the current queue",
    Help: "!pause [url]",
    Function: pause,
}

func pause(s *discordgo.Session, m *discordgo.MessageCreate) {
    guild, ok := GuildMap[m.GuildID]
    if !ok {
        returnMessage := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    if guild.CurrentStream == nil {
        returnMessage := fmt.Sprintf("nothing currently playing")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    currentStream := guild.CurrentStream;

    guild.Paused = !guild.Paused
    currentStream.SetPaused(guild.Paused)

    var message string
    if guild.Paused { message = "paused" } else { message = "unpaused" }
    s.ChannelMessageSend(m.ChannelID, message)
}
