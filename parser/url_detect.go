package parser

type DetectCommand func(string) (UrlType, error)
var DetectFunctions [2]DetectCommand

func initDetectFunctions() {
    DetectFunctions[0] = detectNormal
    DetectFunctions[1] = detectShort
}

func detectNormal(url string) (UrlType, error) {
    return Normal, nil
}

func detectShort(url string) (UrlType, error) {
    return Shorts, nil
}
