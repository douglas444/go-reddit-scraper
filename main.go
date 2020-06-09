package main

import (
    "bufio"
    "os"
    "fmt"
	"strconv"
	"strings"
    "github.com/douglas444/go-reddit-scraper/reddit"
)

func main() {

    reader := bufio.NewReader(os.Stdin);

    fmt.Print("Enter query: ");
    query, _ := reader.ReadString('\n');
    query = strings.TrimSuffix(query, "\n");

    fmt.Print("Enter sort [relevance, hot, top, new, comments]: ");
    sort, _ := reader.ReadString('\n');
    sort = strings.TrimSuffix(sort, "\n");

    fmt.Print("Enter limit: ");
    limitStr, _ := reader.ReadString('\n');
    limitStr = strings.TrimSuffix(limitStr, "\n");
    limit, err := strconv.Atoi(limitStr);

    if err != nil {
        fmt.Println(err);
    }

    posts, err := reddit.Search(query, sort, limit);

    if err != nil {
        fmt.Println(err);
    } else {
        fmt.Println();
        for _, post := range posts {
            fmt.Printf("Upvotes: %d\n", post.Ups);
            fmt.Printf("Downvotes: %d\n", post.Downs);
            fmt.Printf("Comments number: %d\n", post.NumComments);
            fmt.Printf("Media URL: %s\n", post.Url);
            fmt.Printf("Subreddit: %s\n", post.Subreddit);
            fmt.Printf("Title: %s\n\n", post.Title);
        }
    }
}
