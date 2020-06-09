package reddit

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func Search(keyword string) (string, error) {

    client := &http.Client{}
    url := fmt.Sprintf("http://www.reddit.com/search.json?q=%s", keyword);
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err;
    }

    req.Header.Set("User-agent", "douglas444")
    resp, err := client.Do(req)

    if err != nil {
        return "", err;
    }

    defer resp.Body.Close();

    body, err := ioutil.ReadAll(resp.Body);
    if err != nil {
        return "", err;
    }

    return string(body), nil;
}
