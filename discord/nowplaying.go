package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var NowPlaying = Command {
    Name: "np",
    Description: "play audio from youtube videos",
    Help: "!play [url]",
    Function: np,
}

func np(s *discordgo.Session, m *discordgo.MessageCreate) {
    guild, ok := GuildMap[m.GuildID]
    if !ok {
        err := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    if guild.CurrentStream == nil {
        err := fmt.Sprintf("nothing current playing")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    currentSong := guild.Queue[0];
    currentPos := guild.CurrentStream.PlaybackPosition()

    message := fmt.Sprintf("current song: %s\nposition: %s", currentSong.title, currentPos)
    s.ChannelMessageSend(m.ChannelID, message)
}
