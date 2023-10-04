package discord

import (
	"errors"
	"fmt"

	"github.com/jonas747/dca"
)

type Guild struct {
    Id string
    Queue []*Song
    CurrentStream *dca.StreamingSession 
    CurrentSong *Song
    Paused bool
}

func GuildNew(id string) Guild {
    guild := Guild {
        Id: id,
        Queue: make([]*Song, 0),
        CurrentStream: nil,
    }

    return guild
}

func RegisterGuild(id string) error {
    guild := GuildNew(id)
    GuildMap[id] = &guild
    _, ok := GuildMap[id]
    if !ok {
        err := fmt.Sprintf("unable to create new guild")
        return errors.New(err)
    }

    return nil
}
