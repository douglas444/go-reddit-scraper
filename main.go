package main

import (
    "fmt"
    "github.com/douglas444/go-reddit-scraper/reddit"
)

type Job struct {
    Query string
    WindowSize int
    SortBy string
    LastId string
}

func worker(jobs chan Job, results chan Job) {

    for job := range jobs {

        posts, err := reddit.Search(job.Query, job.SortBy, job.WindowSize);

        if err != nil {
            fmt.Println(err);
        }

        cutPoint := len(posts) - 1;
        for i, post := range posts {
            if post.Id == job.LastId {
                cutPoint = i - 1;
                break;
            }
        }

        for i := cutPoint; i >= 0; i-- {
            fmt.Printf("Job query: %s\n", job.Query)
            fmt.Printf("   Id: %s\n", posts[i].Id);
            fmt.Printf("   Upvotes: %d\n", posts[i].Ups);
            fmt.Printf("   Downvotes: %d\n", posts[i].Downs);
            fmt.Printf("   Comments number: %d\n", posts[i].NumComments);
            fmt.Printf("   Subreddit: %s\n", posts[i].Subreddit);
            fmt.Printf("   Title: %s\n\n", posts[i].Title);
        } 

        if len(posts) > 0 && posts[0].Id != job.LastId {
            job.LastId = posts[0].Id;
        }

        results <- job;
    }
}

func main() {

    queries := [3]string{"java", "golang", "javascript"};
    workerPollSize := 2;

    jobs := make(chan Job, len(queries) + 1);
    results := make(chan Job, len(queries));

    for i := 0; i < workerPollSize; i++ {
        go worker(jobs, results);
    } 

    for _, query := range queries {
        jobs <- Job{query, 3, "new", ""};
    }

    for {
        result := <- results;
        jobs <- result;
    }

}


