package reddit

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type Post struct {
    Id string `json:"id"`
    Subreddit string `json:"subreddit"`
    SelfText string `json:"selftext"`
    Downs int `json:"downs"`
    Ups int `json:"ups"`
    NumComments int `json:"num_comments"`
    ThumbnailWidth int `json:"thumbnail_width"`
    ThumbnailHeight int `json:"thumbnail_height"`
    Thumbnail string `json:"thumbnail"`
    Title string `json:"title"`
    Author string `json:"author"`
    CreatedUTC float64 `json:"created_utc"`
    Permalink string `json:"permalink"`
    Url string `json:"url"`
    IsVideo bool `json:"is_video"`
}

type postListing struct {
    Kind string `json:"kind"`
    Data struct {
        Modhash  string `json:"modhash"`
        Children []struct {
            Kind string    `json:"kind"`
            Data Post `json:"data"`
        } `json:"children"`
    } `json:"data"`
}



func Search(query string, sort string, limit int) ([]Post, error) {

    client := &http.Client{}
    urlStr := fmt.Sprintf("https://www.reddit.com/search.json?q=%s&sort=%s&limit=%d", url.QueryEscape(query), sort, limit);

    req, err := http.NewRequest("GET", urlStr, nil)
    if err != nil {
        return nil, err;
    }

    req.Header.Set("User-agent", "douglas444")
    resp, err := client.Do(req)

    if err != nil {
        return nil, err;
    }

    defer resp.Body.Close();

    var result postListing;
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return nil, err;
    }

    var posts []Post;
    for _, post := range result.Data.Children {
        posts = append(posts, post.Data);
    }

    return posts, nil;

}
