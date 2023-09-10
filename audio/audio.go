package audio

import (
	//"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetVideoStreamURL(url string) (string, error) {
    videoSplit := strings.SplitAfter(url, "v=")
    if len(videoSplit) < 1 {
        err := fmt.Sprintf("invalid url")
        return "", errors.New(err)
    }

    videoId := videoSplit[1]

    client := youtube.Client{};
    video, err := client.GetVideo(videoId)
    if err != nil {
        return "", errors.New("unable to get video from id");
    }

    formats := video.Formats
    url, urlErr := client.GetStreamURL(video, &formats[0])
    if urlErr != nil {
        return "", errors.New("unable to create audio stream");
    }

    return url, nil
}

func GetYTVideoInfo(url string) (string, error) {
    videoSplit := strings.SplitAfter(url, "v=")
    if len(videoSplit) < 2 {
        err := fmt.Sprintf("invalid url")
        return "", errors.New(err)
    }

    videoId := videoSplit[1]

    client := youtube.Client{};
    video, err := client.GetVideo(videoId)
    if err != nil {
        return "", errors.New("unable to get video from id");
    }

    return video.Title, nil
}

func GetVideoStream(videoId string) (*os.File, error) {
    client := youtube.Client{}
    video, err := client.GetVideo(videoId)
    if err != nil {
        return nil, errors.New("unable to get video from id");
    }

    formats := video.Formats.WithAudioChannels()
    url, urlErr := client.GetStreamURL(video, &formats[0])
    if urlErr != nil {
        return nil, errors.New("unable to create audio stream");
    }

    fmt.Println(url);

    stream, _, streamErr := client.GetStream(video, &formats[0])
    if streamErr != nil {
        err := fmt.Sprintf("temp err")
        return nil, errors.New(err)
    }

    var tempBuff bytes.Buffer
    var buff bytes.Buffer 
    reader, writer, _ := os.Pipe()

    //remove this if else statement later
    if video.Duration > time.Second * 1 {
        go func() {
            defer writer.Close()
            tempBuff.ReadFrom(stream)
            io.Copy(writer, &tempBuff)
            fmt.Printf("done copying\n")
            fmt.Printf("lengthL: %d\n", buff.Len())
            //f, _ := os.Create("out.mp4")
            //defer f.Close()
            //io.Copy(f, &buff)
        }()
        //time.Sleep(1000*time.Millisecond)
        //go func() {
        //    written, _ := io.Copy(writer, &tempBuff)
        //    bytesWritten = written
        //}()

        return reader, nil
    } else {
        io.Copy(&buff, stream)

        go func() {
            defer writer.Close()
            io.Copy(writer, &buff)
            fmt.Printf("done copying\n")
            fmt.Printf("lengthL: %d\n", buff.Len())
            f, _ := os.Create("out.mp4")
            defer f.Close()
            io.Copy(f, &buff)
        }()

        return reader, nil
    }
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
        err := fmt.Sprintf("unable to use ffmpeg: %s", ffmpegErr.Error())
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
    videoSplit := strings.SplitAfter(url, "v=")
    if len(videoSplit) < 1 {
        err := fmt.Sprintf("invalid url")
        return nil, errors.New(err)
    }

    videoId := videoSplit[1]

    buffReader, streamErr := GetVideoStream(videoId)
    if streamErr != nil {
        err := fmt.Sprintf("unable to get video stream: %s", streamErr.Error())
        return nil, errors.New(err)
    }

    time.Sleep(250 * time.Millisecond)

    go func() {
        encodingErr := EncodeMp4ToMp3(buffReader, writer)
        if encodingErr != nil {
            err := fmt.Sprintf("unable to encode video: %s", encodingErr.Error())
            fmt.Printf("\nissue in encoding go func: %s", err)
            f, _ := os.Create("temp2.mp4")
            defer f.Close()
            io.Copy(f, reader)
            return
            //return buffer, errors.New(err)
        }

        fmt.Printf("done encoding to mp3")
        f, _ := os.Create("temp.mp3")
        defer f.Close()
        io.Copy(f, reader)
    }()

    time.Sleep(500 * time.Millisecond)

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
        err := fmt.Sprintf("unable to get encoded audio: %s", audioBuffErr.Error())
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
        err := fmt.Sprintf("encoding error: %s", encodingErr.Error())
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
