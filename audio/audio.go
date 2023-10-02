package audio

import (
    "main/parser"
	//"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetVideoStreamURL(input string) (string, error) {
   // videoId, parseErr := parser.ParseYTUrl(url)
   // if parseErr != nil {
   //     err := fmt.Sprintf("unable to parse url: %s\n", parseErr.Error())
   //     return "", errors.New(err)
   // }

    id, _ := parser.SetId(input);

    client := youtube.Client{};
    video, err := client.GetVideo(id)
    if err != nil {
        err := fmt.Sprintf("unable to get video from id %s\n", id)
        return "", errors.New(err);
    }

    formats := video.Formats
    url, urlErr := client.GetStreamURL(video, &formats[0])
    if urlErr != nil {
        return "", errors.New("unable to create audio stream\n");
    }

    return url, nil
}

func GetYTVideoInfo(url string) (string, error) {
    videoId, parseErr := parser.ParseYTUrl(url)
    if parseErr != nil {
        err := fmt.Sprintf("unable to parse url: %s\n", parseErr.Error())
        return "", errors.New(err)
    }

    client := youtube.Client{};
    video, err := client.GetVideo(videoId)
    if err != nil {
        err := fmt.Sprintf("unable to get video from id %s\n", videoId)
        return "", errors.New(err);
    }


    return video.Title, nil
}

func GetVideoStream(videoId string) (*os.File, error) {
    client := youtube.Client{}
    video, err := client.GetVideo(videoId)
    if err != nil {
        return nil, errors.New("unable to get video from id\n");
    }

    formats := video.Formats.WithAudioChannels()
    url, urlErr := client.GetStreamURL(video, &formats[0])
    if urlErr != nil {
        return nil, errors.New("unable to create audio stream\n");
    }

    fmt.Println(url);

    stream, _, streamErr := client.GetStream(video, &formats[0])
    if streamErr != nil {
        err := fmt.Sprintf("unable to get yt stream\n")
        return nil, errors.New(err)
    }

    var buff bytes.Buffer 
    reader, writer, _ := os.Pipe()

    go func() {
        defer writer.Close()
        buff.ReadFrom(stream)
        io.Copy(writer, &buff)
    }()

    time.Sleep(500 * time.Millisecond)

    return reader, nil
}

func EncodeMp4ToMp3(stream *os.File, audioBuff *io.PipeWriter) error {
    ffmpegErr := ffmpeg.Input(
            "pipe:0",  
        ).
        Output("pipe:1",
                         ffmpeg.KwArgs{"b:a": "84K"},
                         ffmpeg.KwArgs{"vn": ""},
                         ffmpeg.KwArgs{"c:a": "libmp3lame"},
                         ffmpeg.KwArgs{"f": "mp3"}).
                OverWriteOutput().
                WithOutput(audioBuff).
                ErrorToStdOut().
                WithInput(stream).
                Run()


    if ffmpegErr != nil {
        err := fmt.Sprintf("unable to use ffmpeg: %s\n", ffmpegErr.Error())
        fmt.Printf("ffmpeg error, closing")
        audioBuff.Close()
        return errors.New(err)
    }

    fmt.Printf("done with ffmpeg encoding\n")
    audioBuff.Close()

    return nil
}

func GetYTAudioBuffer(url string) (*io.PipeReader, error) {
    //do some more error handling for the video url thing
    //no parsing for short urls
    reader, writer := io.Pipe()

    videoId, parseErr := parser.ParseYTUrl(url)
    if parseErr != nil {
        err := fmt.Sprintf("unable to parse url: %s\n", parseErr.Error())
        return nil, errors.New(err)
    }

    buffReader, streamErr := GetVideoStream(videoId)
    if streamErr != nil {
        err := fmt.Sprintf("unable to get video stream: %s\n", streamErr.Error())
        return nil, errors.New(err)
    }

    time.Sleep(500 * time.Millisecond)

    go func() {
        encodingErr := EncodeMp4ToMp3(buffReader, writer)
        if encodingErr != nil {
            err := fmt.Sprintf("unable to encode video: %s\n", encodingErr.Error())
            fmt.Printf("\nissue in encoding go func: %s", err)
            //return buffer, errors.New(err)
        }
    }()

    return reader, nil;
}

func SaveAudioBuffer(buff *bytes.Buffer) {
    f, _ := os.Create("out.ogg")
    f.Write(buff.Bytes())
    f.Close()
}

func TestEncoder() {
    url := "https://www.youtube.com/watch?v=c8LSRGJO5Mk"

    audioBuff, audioBuffErr := GetYTAudioBuffer(url)
    if audioBuffErr != nil {
        err := fmt.Sprintf("unable to get encoded audio: %s\n", audioBuffErr.Error())
        fmt.Printf("\n%s", err)
        return
    }

    time.Sleep(250 * time.Millisecond)

    options := dca.StdEncodeOptions
    options.RawOutput = true
    options.Bitrate = 96
    options.Application = "lowdelay"
    options.Volume = 500

    encodingSession, encodingErr := dca.EncodeMem(audioBuff, options)
    if encodingErr != nil {
        err := fmt.Sprintf("encoding error: %s\n", encodingErr.Error())
        fmt.Printf(err)
        return
    }

    time.Sleep(2000 * time.Millisecond)

    fmt.Printf("starting ticker\n")
    ticker := time.NewTicker(time.Second)
    for {
        select {
        case <-ticker.C: {
            stats := encodingSession.Stats()
            err := encodingSession.Error()
            if err != nil {
                fmt.Printf("error while encoding: %s\n", encodingSession.FFMPEGMessages())
            }

            fmt.Printf("transcode status: time: %s\n", stats.Duration)
        }
        }
    }
}
