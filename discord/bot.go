package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
)

var GuildMap map[string]*Guild

var (
    Token string
    CmdId byte
)

func init() {
    GuildMap = make(map[string]*Guild)
    initCommands()
}

func Start() {
    envErr := godotenv.Load()
    if envErr != nil {
        fmt.Println("unable to open env file, closing")
        os.Exit(1);
    }

    Token = os.Getenv("Token")
    CmdId = os.Getenv("CmdId")[0]

    dg, connErr := discordgo.New("Bot " + Token)
    if connErr != nil {
        fmt.Printf("unable to connect bot to discord servers: %s", connErr.Error())
        os.Exit(1)
    }

    dg.AddHandler(messageCreate)

    dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
    
    openErr := dg.Open()
    if openErr != nil {
        fmt.Println("unable to open connection to discord")
        return
    }

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID { return }
    if m.Content[0] == CmdId { handleCommand(s, m) }
}
