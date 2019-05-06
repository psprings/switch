package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// CurrentTimeString :
func CurrentTimeString() string {
	now := time.Now().UTC()
	return now.String()
}

// GetHookURLFromRequest :
func GetHookURLFromRequest(r *http.Request) string {
	urlParam, found := r.URL.Query()["url"]
	println(urlParam)
	if found {
		return urlParam[0]
	}
	postToURL := r.URL.Path
	postToURL = strings.TrimPrefix(postToURL, "/")
	return "http://" + postToURL
}

func isURLServerError(res *http.Response) bool {
	r, _ := regexp.Compile(`^5\d\d$`)
	return r.MatchString(string(res.StatusCode))
}

// TestEndpoint :
func TestEndpoint(testURL string) error {
	res, err := http.Get(testURL)
	if err != nil {
		return err
	}
	if isURLServerError(res) {
		return errors.New("5xx error for URL " + testURL)
	}
	return nil
}

// ForwardPostRequest :
func ForwardPostRequest(forwardToURL string, r *http.Request) {
	client := &http.Client{}
	u, err := url.Parse(forwardToURL)
	if err != nil {
		log.Fatal(err)
	}
	r.URL = u
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	fmt.Println("response headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body:", string(body))
}

// EnsureBodyContent :
func EnsureBodyContent(body string) string {
	if body == "" {
		return body + "Placeholder."
	}
	return body
}
