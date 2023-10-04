package discord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Jump = Command {
    Name: "jump",
    Aliases: []string{"jump", "jmp"},
    Description: "remove an item from the queue",
    Help: "!remove [queue position]",
    Function: jump,
}

func jump(s *discordgo.Session, m *discordgo.MessageCreate) {
    argumentSplit := strings.Split(m.Content, " ")
    if len(argumentSplit) <= 1 {
        s.ChannelMessageSend(m.ChannelID, "not enough arguments")
        return 
    }

    guild, ok := GuildMap[m.GuildID]
    if !ok {
        returnMessage := fmt.Sprintf("guild is unregistered, please play at least one song")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    queueLen := len(guild.Queue)

    if queueLen == 0 {
        returnMessage := fmt.Sprintf("queue is empty")
        s.ChannelMessageSend(m.ChannelID, returnMessage)
        return
    }

    requestedPosition, convErr := strconv.Atoi(argumentSplit[1])
    if convErr != nil {
        err := fmt.Sprintf("input cannot be converted to integer")
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    if requestedPosition > queueLen || requestedPosition == 0 {
        err := "invalid index"
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    guild.Queue = guild.Queue[requestedPosition-1:]
    guild.CurrentStream.Kill()
}
