package main

import (
    "fmt"
    "github.com/douglas444/go-reddit-scraper/reddit"
)

type Job struct {
    Query string
    LastId string
}

func worker(jobs chan Job, results chan Job) {

    for job := range jobs {

        posts, err := reddit.Search(job.Query, "new", 1);

        if err != nil {
            fmt.Println(err);
        }

        if len(posts) > 0 && posts[0].Id != job.LastId {

            job.LastId = posts[0].Id;

            fmt.Printf("Job query: %s\n", job.Query)
            fmt.Printf("   Id: %s\n", posts[0].Id);
            fmt.Printf("   Upvotes: %d\n", posts[0].Ups);
            fmt.Printf("   Downvotes: %d\n", posts[0].Downs);
            fmt.Printf("   Comments number: %d\n", posts[0].NumComments);
            fmt.Printf("   Subreddit: %s\n", posts[0].Subreddit);
            fmt.Printf("   Title: %s\n\n", posts[0].Title);
        }

        results <- job;
    }
}

func main() {

    queries := [3]string{"trump", "bolsonaro", "nicolás maduro"};
    workerPollSize := 2;

    jobs := make(chan Job, len(queries) + 1);
    results := make(chan Job, len(queries));

    for i := 0; i < workerPollSize; i++ {
        go worker(jobs, results);
    } 

    for _, query := range queries {
        jobs <- Job{query, ""};
    }

    for {
        result := <- results;
        jobs <- result;
    }

}


