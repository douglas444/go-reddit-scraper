package reddit

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func Search(keyword string) (string, error) {

    url := fmt.Sprintf("http://www.reddit.com/search.json?q=%s", keyword);
	resp, err := http.Get(url);

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
