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

	"github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetYTVideoInfo(url string) (*YTInfo, error) {
    videoId, idErr := parser.SetId(url);
    if idErr != nil {
        err := fmt.Sprintf("unable to get video id: %s", idErr.Error())
        return nil, errors.New(err)
    }

    client := youtube.Client{};
    video, err := client.GetVideo(videoId)
    if err != nil {
        err := fmt.Sprintf("unable to get video from id while getting yt info%s\n", videoId)
        return nil, errors.New(err);
    }

    v := YTInfo {
        Title: video.Title,
        Length: video.Duration.String(),
        Id: video.ID,
    }

    return &v, nil
}

func GetVideoStream(videoId string) (*os.File, error) {
    client := youtube.Client{}
    video, err := client.GetVideo(videoId)
    if err != nil {
        return nil, errors.New("unable to get video from id while getting stream\n");
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

    //videoId, parseErr := parser.ParseYTUrl(url)
    //if parseErr != nil {
    //    err := fmt.Sprintf("unable to parse url: %s\n", parseErr.Error())
    //    return nil, errors.New(err)
    //}

    videoId, _ := parser.SetId(url);

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
