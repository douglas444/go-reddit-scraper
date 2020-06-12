package main

import (
    "fmt"
    "strconv"
    "strings"
    "net/http"
    "github.com/douglas444/go-reddit-scraper/reddit"
)

type Job struct {
    Id int
    Query string
    WindowSize int
    SortBy string
    LastId string
    IsActive bool
}

func scrape(query string, sortBy string, windowSize int, lastId string) string {

    posts, err := reddit.Search(query, sortBy, windowSize);

    if err != nil {
        fmt.Println(err);
    }

    cutPoint := len(posts) - 1;
    for i, post := range posts {
        if post.Id == lastId {
            cutPoint = i - 1;
            break;
        }
    }

    for i := cutPoint; i >= 0; i-- {
        fmt.Println(query, "|", posts[i].Title);
    }

    if len(posts) > 0 && posts[0].Id != lastId {
        return posts[0].Id;
    } else {
        return lastId;
    }
}

func worker(jobs chan Job) {

    for job := range jobs {
 
        if job.IsActive {
            job.LastId = scrape(job.Query, job.SortBy, job.WindowSize, job.LastId);
        }

        jobs <- job;
    }
}

func serverStart(exit chan bool, jobById map[int]Job) {
    
    http.HandleFunc("/exit", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("[SERVER LOG] exiting...");
        exit <- true;
    });

    http.HandleFunc("/deactivate/", func(w http.ResponseWriter, req *http.Request) {

        jobId, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/deactivate/"));

        if err != nil {
		    w.WriteHeader(400);
            return;
        } else if _, isPresent := jobById[jobId]; !isPresent {
		    w.WriteHeader(400);
		    return;
        } else {
            job := jobById[jobId];
            job.IsActive = false;
            fmt.Println("[SERVER LOG] job", job.Id, "deactivated");
        }
    });

    http.HandleFunc("/activate/", func(w http.ResponseWriter, req *http.Request) {

        jobId, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/activate/"));

        if err != nil {
		    w.WriteHeader(400);
		    return;
        } else if _, isPresent := jobById[jobId]; !isPresent {
		    w.WriteHeader(400);
		    return;
        } else {
            job := jobById[jobId];
            job.IsActive = true;
            fmt.Println("[SERVER LOG] job", job.Id, "activated");
        }
    });

    http.ListenAndServe(":8080", nil);

}

func main() {

    queries := [3]string{"bolsonaro", "trump", "nicolás maduro"};
    workerPollSize := 2;

    jobs := make(chan Job, len(queries) + 1);

    for i := 0; i < workerPollSize; i++ {
        go worker(jobs);
    }

    jobById := make(map[int]Job);

    for id, query := range queries {
        job := Job{id, query, 3, "new", "", true};
        jobById[id] = job;
        jobs <- job;
    }

    exit := make(chan bool);

    go serverStart(exit, jobById);

    <- exit;

}


