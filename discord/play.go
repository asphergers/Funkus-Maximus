package discord

import (
	_ "bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"main/audio"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

var Play = Command {
    Name: "play",
    Description: "play audio from youtube videos",
    Help: "!play [url]",
    Function: play,
}

func play(s *discordgo.Session, m *discordgo.MessageCreate) {
    argumentSplit := strings.SplitN(m.Content, " ", 2)
    if len(argumentSplit) < 1 {
        s.ChannelMessageSend(m.ChannelID, "not enough arguments")
        return 
    }

    url := argumentSplit[1]

    currentUserVC, userChannelErr := getCurrentUserChannel(s, m)
    if userChannelErr != nil {
        err := fmt.Sprintf("you are not in a vc: %s\n", userChannelErr.Error())
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    //guild stuff
    guild, ok := GuildMap[m.GuildID];
    if !ok {
        err := RegisterGuild(m.GuildID)
        if err != nil {
            err := fmt.Sprintf("unable to register guild: %s\n", err.Error())
            s.ChannelMessageSend(m.ChannelID, err)
            return
        }

        guild = GuildMap[m.GuildID]
    }

    fmt.Println(guild.Id)
    fmt.Println(guild.CurrentStream)

    if guild.CurrentStream != nil {
        title, err := audio.GetYTVideoInfo(url);
        if err != nil {
            returnErr := fmt.Sprintf("unable to get yt video info: %s", err.Error())
            s.ChannelMessageSend(m.ChannelID, returnErr)
            return
        }

        song := Song {
            url: url,
            title: title,
            buff: nil,
        }

        AddSongToQueue(guild, &song)

        addMessage := fmt.Sprintf("added song to queue: %s", song.title)
        s.ChannelMessageSend(m.ChannelID, addMessage)
        return
    }

    title, infoErr := audio.GetYTVideoInfo(url)
    if infoErr != nil {
        err := fmt.Sprintf("unable to get video info: %s\n", infoErr.Error())
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    song := Song {
        title: title,
        url: url,
        buff: nil,
    }
    
    AddSongToQueue(guild, &song)

    vc, vcJoinErr := s.ChannelVoiceJoin(m.GuildID, currentUserVC, false, true)
    if vcJoinErr != nil {
        err := fmt.Sprintf("unable to join channel: %s\n", vcJoinErr.Error())
        s.ChannelMessageSend(m.ChannelID, err)
        vc.Disconnect()
        return
    }

    defer func() {
        vc.Speaking(false)
        vc.Disconnect()
        guild.CurrentStream = nil
        guild.CurrentSong = nil
    }()

    for (len(guild.Queue) > 0) {
        song := guild.Queue[0]
        var nextSong *Song
        if len(guild.Queue) >= 2 { nextSong = guild.Queue[1] }

        message := fmt.Sprintf("now playing: %s", song.title)
        s.ChannelMessageSend(m.ChannelID, message)

        go func() {
            if nextSong == nil { return }
            audioBuff, buffErr := audio.GetYTAudioBuffer(nextSong.url)
            if buffErr != nil { return }
            nextSong.buff = audioBuff
        }()

        guild.CurrentSong = song
        guild.Queue = guild.Queue[1:]

        PlaySong(s, m, song, vc, guild)
    }
}

func PlaySong(s *discordgo.Session, m *discordgo.MessageCreate, song *Song, 
                vc *discordgo.VoiceConnection, guild *Guild) {


    if song.buff == nil {
        audioBuff, audioBuffErr := audio.GetYTAudioBuffer(song.url)
        if audioBuffErr != nil {
            err := fmt.Sprintf("unable to get encoded audio: %s\n", audioBuffErr.Error())
            s.ChannelMessageSend(m.ChannelID, err)
            return
        }

        song.buff = audioBuff
    }

    options := dca.StdEncodeOptions
    options.RawOutput = true
    options.Bitrate = 96
    options.Application = "lowdelay"
    options.Volume = 500

    encodingSession, encodingErr := dca.EncodeMem(song.buff, options)
    if encodingErr != nil {
        err := fmt.Sprintf("encoding error: %s\n", encodingErr.Error())
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    //time.Sleep(250 * time.Millisecond)

    done := make(chan error)

    defer func() {
        encodingSession.Cleanup()
        fmt.Println("done")
    }()

    speakErr := vc.Speaking(true)
    if speakErr != nil {
        err := fmt.Sprintf("unable to fucking speak: %s\n", speakErr.Error())
        s.ChannelMessageSend(m.ChannelID, err)
        return
    }

    stream := dca.NewStream(encodingSession, vc, done)
    ticker := time.NewTicker(time.Second)

    guild.CurrentStream = stream

    for {
        select {
        case err := <- done: {
            if err != nil && err != io.EOF {
                fmt.Printf("error while streaming: %s", err.Error())
            }

            return
        }

        case <-ticker.C: {
            stats := encodingSession.Stats()
            pos := stream.PlaybackPosition()

            fmt.Printf("playback %s, transcode status: time: %s\n", pos, stats.Duration)
        }
        }
    }
}

func getCurrentUserChannel(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
    guild, err := s.State.Guild(m.GuildID)
    if err != nil { return "", errors.New("unable to find your current guild(?)") }

    for _, vs := range guild.VoiceStates {
        if vs.UserID == m.Author.ID {
            return vs.ChannelID, nil
        }
    }

    return "", nil
}

func AddSongToQueue(guild *Guild, song *Song) error {
    guild.Queue = append(guild.Queue, song)

    return nil
}
