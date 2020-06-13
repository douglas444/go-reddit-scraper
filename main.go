package main

import (
    "fmt"
    "strconv"
    "strings"
    "net/http"
    "github.com/douglas444/go-reddit-scraper/reddit"
)

type RequestType int

const(
    Exit RequestType = 1
    ActivateJob RequestType = 2
    DeactivateJob RequestType = 3
)

type Request struct {
    Type RequestType
    JobId int
}

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
        fmt.Println(query, "|", posts[i].Title, "\n");
    }

    if len(posts) > 0 && posts[0].Id != lastId {
        return posts[0].Id;
    } else {
        return lastId;
    }
}

func worker(jobs chan *Job) {

    for job := range jobs {
 
        if job.IsActive {
            job.LastId = scrape(job.Query, job.SortBy, job.WindowSize, job.LastId);
            jobs <- job;
        }

    }
}

func requestProcessor(requests chan Request, exit chan bool, jobs chan *Job, jobById map[int]*Job) {

    for request := range requests {
 
        switch request.Type {

        case Exit:

            exit <- true;
            fmt.Println("exiting");

        case ActivateJob:

            if _, isPresent := jobById[request.JobId]; !isPresent {
                fmt.Println("ignoring activate request for invalid job id");
            } else if job, _ := jobById[request.JobId]; job.IsActive {
                fmt.Println("ignoring activate request for already active job");
            } else {
                jobById[request.JobId].IsActive = true;
                jobs <- jobById[request.JobId];
                fmt.Println("activating job", request.JobId);
            }

        case DeactivateJob:

            if _, isPresent := jobById[request.JobId]; !isPresent {
                fmt.Println("ignoring request for invalid job id");
            } else if job, _ := jobById[request.JobId]; !job.IsActive {
                fmt.Println("ignoring deactivate request for already deactive job");
            } else {
                jobById[request.JobId].IsActive = false;     
                fmt.Println("deactivating job", request.JobId);
            }
        }
    }
}

func serverStart(requests chan Request) {
    
    http.HandleFunc("/exit", func(w http.ResponseWriter, req *http.Request) {

        request := Request{Exit, -1};
        select {
        case requests <- request:
            break;
        default:
            w.WriteHeader(503);
        }

    });

    http.HandleFunc("/deactivate/", func(w http.ResponseWriter, req *http.Request) {

        jobId, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/deactivate/"));

        if err != nil {
            w.WriteHeader(400);
            return;
        }

        request := Request{DeactivateJob, jobId};
        select {
        case requests <- request:
            break;
        default:
            w.WriteHeader(503);
        }

    });

    http.HandleFunc("/activate/", func(w http.ResponseWriter, req *http.Request) {

        jobId, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/activate/"));

        if err != nil {
            w.WriteHeader(400);
            return;
        }

        request := Request{ActivateJob, jobId};
        select {
        case requests <- request:
            break;
        default:
            w.WriteHeader(503);
        }

    });

    http.ListenAndServe(":8080", nil);

}

func main() {

    queries := [3]string{"bolsonaro", "trump", "nicolás maduro"};
    workerPollSize := 2;
    requestsChannelSize := 3;

    jobs := make(chan *Job, len(queries) + 1);

    for i := 0; i < workerPollSize; i++ {
        go worker(jobs);
    }

    jobById := make(map[int]*Job);

    for id, query := range queries {
        job := Job{id, query, 3, "new", "", true};
        jobById[id] = &job;
        jobs <- &job;
    }

    requests := make(chan Request, requestsChannelSize);
    exit := make(chan bool);

    go requestProcessor(requests, exit, jobs, jobById);
    go serverStart(requests);

    <- exit;

}


