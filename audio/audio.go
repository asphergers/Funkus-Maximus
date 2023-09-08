package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hraban/opus"
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

func GetVideoStream(videoId string) (io.ReadCloser, error) {
    client := youtube.Client{};
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

    result, getErr := http.Get(url);
    if getErr != nil {
        return nil, errors.New("unable to get stream from url");
    }

    return result.Body, nil;
}

func EncodeMp4ToMp3(stream io.ReadCloser, audioBuff *io.PipeWriter) error {
    ffmpegErr := ffmpeg.Input(
            "pipe:0",  
        ).
        Output("pipe:1",
                         ffmpeg.KwArgs{"b:a": "96K"},
                         ffmpeg.KwArgs{"vn": ""},
                         ffmpeg.KwArgs{"f": "mp3"}).
                OverWriteOutput().
                WithOutput(audioBuff).
                WithInput(stream).
                Run()


    if ffmpegErr != nil {
        err := fmt.Sprintf("unable to use ffmpeg: %s", ffmpegErr.Error())
        return errors.New(err);
    }

    audioBuff.Close()

    return nil
}

func encodeAudio(stream io.ReadCloser, audioBuff *bytes.Buffer) error {
    ffmpegErr := ffmpeg.Input(
            "pipe:0",  
        ).
        Output("pipe:1",
                         ffmpeg.KwArgs{"c:a": "libopus"},
                         ffmpeg.KwArgs{"b:a": "96K"},
                         ffmpeg.KwArgs{"ac": "2"},
                         ffmpeg.KwArgs{"vn": ""},
                         ffmpeg.KwArgs{"ar": "48000"},
                         ffmpeg.KwArgs{"f": "ogg"}).
                         WithInput(stream).WithOutput(audioBuff).Run()


    if ffmpegErr != nil {
        err := fmt.Sprintf("unable to use ffmpeg: %s", ffmpegErr.Error())
        return errors.New(err);
    }

    return nil
}

func EncodeOpusBuff(outbuff [][]byte, inBuff []int16) error {
    return nil
}

func DecodeOpusBuff(outBuff [][]int16, inBuff *bytes.Buffer) error {
    channels := 2
    s, err := opus.NewStream(inBuff)
    if err != nil {
        err := fmt.Sprintf("unable to open opus stream: %s", err.Error())
        return errors.New(err)
    }

    defer s.Close()
    chunk := make([]int16, 16384)
    for {
        n, err := s.Read(chunk)
        if err == io.EOF {
            return nil
        } else if err != nil {
            err := fmt.Sprintf("unxepceted error reading opus stream: %s", err.Error())
            return errors.New(err)
        }

        outBuff = append(outBuff, chunk[:n*channels])
    }
}

func FormatOpusBuff(outBuff [][]byte, inBuff *bytes.Buffer) error {
    var opusLen int16

    for {
        err := binary.Read(inBuff, binary.LittleEndian, &opusLen)

        if err == io.EOF || err == io.ErrUnexpectedEOF {
            return err
        }

        if err != nil {
            returnErr := fmt.Sprintf("unable to read opus length: %s", err.Error())
            return errors.New(returnErr)
        }

        fmt.Printf("opus len: %d", opusLen)
        buf := make([]byte, opusLen)
        readErr := binary.Read(inBuff, binary.LittleEndian, &buf)
        if readErr != nil {
            err := fmt.Sprintf("unable to read audio data into second part of opus stream: %s", readErr.Error()) 
            return errors.New(err);
        }

        outBuff = append(outBuff, buf)
    }

}

func GetYTAudioBuffer(url string) (*io.PipeReader, error) {
    //do some more error handling for the video url thing
    //no parsing for short urls
    buffReader, buffWriter := io.Pipe()
    videoSplit := strings.SplitAfter(url, "v=")
    if len(videoSplit) < 1 {
        err := fmt.Sprintf("invalid url")
        return buffReader, errors.New(err)
    }

    videoId := videoSplit[1]

    stream, streamErr := GetVideoStream(videoId)
    if streamErr != nil {
        err := fmt.Sprintf("unable to get video stream: %s", streamErr.Error())
        return buffReader, errors.New(err)
    }

    go func() {
        encodingErr := EncodeMp4ToMp3(stream, buffWriter)
        if encodingErr != nil {
            err := fmt.Sprintf("unable to encode video: %s", encodingErr.Error())
            fmt.Printf("\nissue in encoding go func: %s", err)
            //return buffer, errors.New(err)
        }
    }()

    time.Sleep(500 * time.Millisecond)

    return buffReader, nil;
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

func Test() {
    buffer := bytes.NewBuffer(make([]byte, 0))
    videoId := "O8tGBmyYirc" 

    stream, streamErr := GetVideoStream(videoId)
    if streamErr != nil {
        fmt.Printf("unable to get video stream: %s", streamErr.Error())
        return
    }

    encodingErr := encodeAudio(stream, buffer)
    if encodingErr != nil {
        fmt.Printf("unable to encode video: %s", encodingErr.Error())
    }
}
