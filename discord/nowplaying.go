package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var NowPlaying = Command {
    Name: "np",
    Aliases: []string{"nowplaying", "np"},
    Description: "get information about the current track",
    Help: "!nowplaying [url]",
    Function: np,
}

func np(s *discordgo.Session, m *discordgo.MessageCreate) {
    guild, ok := GuildMap[m.GuildID]
    if !ok {
        returnMessage := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    if guild.CurrentStream == nil {
        returnMessage := fmt.Sprintf("nothing current playing")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    currentSong := guild.CurrentSong
    currentPos := guild.CurrentStream.PlaybackPosition()

    message := fmt.Sprintf("current song: %s\nposition: %s\nlength: %s", currentSong.title, currentPos, currentSong.length)
    s.ChannelMessageSend(m.ChannelID, message)
}
