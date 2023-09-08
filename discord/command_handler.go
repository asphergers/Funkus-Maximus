package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const COMMANDS_SIZE = 10;

var commandList [COMMANDS_SIZE]string
var functions [COMMANDS_SIZE]Command

func initCommands() {
    commandList[0] = "ping"
    functions[0] = Ping

    commandList[1] = "test"
    functions[1] = Test

    commandList[2] = "play"
    functions[2] = Play

    commandList[3] = "skip"
    functions[3] = Skip

    commandList[4] = "np"
    functions[4] = NowPlaying

    commandList[5] = "q"
    functions[5] = Queue
}

func handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    spaceIndex := strings.Index(m.Message.Content, " ");
    found := false;

    var command string
    if spaceIndex == -1 { command = m.Content[1:] } else { command = m.Content[1:spaceIndex] }

    for i := 0; i < COMMANDS_SIZE; i++ {
        if command == commandList[i] {
            found = true;
            functions[i].Function(s, m);
        }
    }
    
    if !found { s.ChannelMessageSend(m.ChannelID, "invalid command!") }
}
