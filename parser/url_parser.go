package parser

import (
	"fmt"
	"strings"
    "errors"
)

type UrlType int
const (
    Normal UrlType = 0
    Shorts UrlType = 1
    Shortened UrlType = 2
)

func InitParser() {
    initDetectFunctions()
}

//this does not work
//i don't knwo why it doesnt work
//i don't want to know why is doesn't work
//this should work
//return 0 for now so I can test
func getUrlType(url string) UrlType {
    for _, function := range DetectFunctions {
        t, err := function(url)
        if err != nil { return t }
    }

    return 0
}

func ParseYTUrl(url string) (string, error) {
    var id string
    var err error

    urlType := getUrlType(url);
    if urlType == -1 {
        err := fmt.Sprintf("invalid url type\n")
        return "", errors.New(err)
    }

    switch urlType {
    case Normal: { id, err = parseNormalUrl(url) }
    case Shorts: { id, err = parseShortsUrl(url) }

    default: { id,err = parseNormalUrl(url) }
    }

    if err != nil { return "", err }

    return id, nil
}

func parseNormalUrl(url string) (string, error) {
    var urlId string

    if strings.Contains(url, "?si=") {
        leftDelim := "?v="
        rightDelim := "?si="

        id, err := getInbetween(url, leftDelim, rightDelim)
        if err != nil {
            err := fmt.Sprintf("invalid normal url: %s\n", err.Error())
            return "", errors.New(err)
        }

        urlId = id

    } else {
        videoSplit := strings.SplitAfter(url, "v=")
        if len(videoSplit) < 2 {
            err := fmt.Sprintf("invalid url\n")
            return "", errors.New(err)
        }

        urlId = videoSplit[1]
    }

    if urlId == "" {
        err := fmt.Sprintf("invalid url\n")
        error := errors.New(err)
        return "", error
    }

    return urlId, nil
}

func parseShortsUrl(url string) (string, error) {
    var urlId string

    if strings.Contains(url, "?si=") {
        leftDelim := "/shorts/"
        rightDelim := "?si="

        id, err := getInbetween(url, leftDelim, rightDelim)
        if err != nil {
            err := fmt.Sprintf("invalid shorts url: %s\n", err.Error())
            return "", errors.New(err)
        }

        urlId = id

    } else {
        urlSplit := strings.Split(url, "/shorts/")
        if len(urlSplit) < 2 {
            err := fmt.Sprintf("invalid shorts url\n")
            return "", errors.New(err)
        }

        urlId = urlSplit[1]
    }

    if urlId == "" {
        err := fmt.Sprintf("invalid shorts url\n")
        error := errors.New(err)
        return "", error
    }

    return urlId, nil
}

func getInbetween(s string, leftDelim string, rightDelim string) (string, error) {
    leftIndex := strings.Index(s, leftDelim)
    if leftIndex == -1 {
        err := fmt.Sprintf("invlaid string no left delimeter\n")
        return "", errors.New(err)
    }

    rightIndex := strings.Index(s, rightDelim)
    if rightIndex == -1 {
        err := fmt.Sprintf("invlaid string no right delimeter\n")
        return "", errors.New(err)
    }

    leftIndex += len(leftDelim)
    
    subStr := s[leftIndex:rightIndex]

    return subStr, nil 
}
