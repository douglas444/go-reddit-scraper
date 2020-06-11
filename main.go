package main

import (
    "bufio"
    "os"
    "fmt"
    "strings"
    "github.com/douglas444/go-reddit-scraper/reddit"
)

func exec(query string, c chan reddit.Post) {

    posts, err := reddit.Search(query, "new", 1);

    if err != nil {
        fmt.Println(err);
    }

    if len(posts) > 0 {
        c <- posts[0];
    }
}

func main() {

    reader := bufio.NewReader(os.Stdin);

    fmt.Print("Enter query: ");
    query, _ := reader.ReadString('\n');
    query = strings.TrimSuffix(query, "\n");

    c := make(chan reddit.Post);
    var lastId string;

    for {
        go exec(query, c);
        select {
            case post := <- c:
                if post.Id != lastId {
                    lastId = post.Id;
                    fmt.Printf("Upvotes: %d\n", post.Ups);
                    fmt.Printf("Downvotes: %d\n", post.Downs);
                    fmt.Printf("Comments number: %d\n", post.NumComments);
                    fmt.Printf("Media URL: %s\n", post.Url);
                    fmt.Printf("Subreddit: %s\n", post.Subreddit);
                    fmt.Printf("Title: %s\n\n", post.Title);
                }
        }
    }
}


