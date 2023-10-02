package parser

import (
    "fmt"
    "errors"

    ytsearch "github.com/lithdew/youtube" 
)

func SearchYoutube(query string) (string, error) {
    results, err := ytsearch.Search(query, 0);
    if err != nil {
        err := fmt.Sprintf("invalid search: %s", err.Error()) 
        return "", errors.New(err)
    }

    result_id := fmt.Sprintf("%q", results.Items[0].ID);

    return result_id, nil
}

func SetId(input string) (string, error) {
    if len(input) < len("https://") {
        tid, err := SearchYoutube(input)
        if err != nil {
            e := fmt.Sprintf("bad query %s\n", input)
            return "", errors.New(e);
        }

        return tid, nil
    }
    
    if input[:8] != "https://" {
        tid, err := SearchYoutube(input)
        if err != nil {
            e := fmt.Sprintf("bad query %s\n", input)
            return "", errors.New(e);
        }

        return tid, nil
    }

    return input, nil;
}
