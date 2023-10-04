package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)


var functions [8]Command

func initCommands() {
    functions[0] = Ping
    functions[1] = Test
    functions[2] = Play
    functions[3] = Skip
    functions[4] = NowPlaying
    functions[5] = Queue
    functions[6] = Pause
    functions[7] = Remove
}

func handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    spaceIndex := strings.Index(m.Message.Content, " ");
    found := false;

    var command string
    if spaceIndex == -1 { command = m.Content[1:] } else { command = m.Content[1:spaceIndex] }

    for i := 0; i < len(functions); i++ {
        for j := 0; j < len(functions[i].Aliases); j++ {
            if command == functions[i].Aliases[j] {
                found = true
                functions[i].Function(s, m)
                break
            }
        }

        if found { break }
    }
    
    if !found { s.ChannelMessageSend(m.ChannelID, "invalid command!") }
}
