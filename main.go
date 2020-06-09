package main

import (
    "fmt"
	"github.com/douglas444/go-reddit-scraper/reddit"
)

func main() {

    body, err := reddit.Search("teste");

    if err != nil {
        fmt.Println(err);
    } else {
        fmt.Println(body);
    }
}
