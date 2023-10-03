package parser

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func ytdl_search_first(query string) (string, string, error) {
    cmd := exec.Command("youtube-dl", "--skip-download", "--get-title", "--get-id", "ytsearch1:" + query);

    var stdOut, stdErr bytes.Buffer
    cmd.Stdout = &stdOut
    cmd.Stderr = &stdErr
    runErr := cmd.Run()
    if runErr != nil {
        err := fmt.Sprintf("unable to query using youtube dl: %s", runErr.Error())
        return "", "", errors.New(err)
    }

    split := strings.Split(stdOut.String(), "\n")
    if len(split) < 2 {
        err := fmt.Sprintf("invalid youtube dl parsed response: %s", stdOut.String())
        return "", "", errors.New(err)
    }

    return split[0], split[1], nil
}

func SearchYoutube(query string) (string, error) {
    _, id, searchErr := ytdl_search_first(query)
    if searchErr != nil {
        err := fmt.Sprintf("unable to search: %s", searchErr.Error())
        return "", errors.New(err)
    }

    return id, nil
}

func SetId(input string) (string, error) {
    if len(input) < len("https://") {
        tid, err := SearchYoutube(input)
        if err != nil {
            e := fmt.Sprintf("bad query: %s\n", err.Error())
            return "", errors.New(e);
        }

        return tid, nil
    }
    
    if input[:8] != "https://" {
        tid, err := SearchYoutube(input)
        if err != nil {
            e := fmt.Sprintf("bad query: %s\n", err.Error())
            return "", errors.New(e);
        }

        return tid, nil
    }

    return input, nil;
}
