package memo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"errors"
	"time"
)


func httpGetBody(url string) (interface{}, error) {
	if ! isValidHttpURL(url) {
		e := errors.New("Invalid URL string specified")
		return nil, e
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func TestMemo_Get(t *testing.T) {
	urllist := []string{
		"https://www.google.com",
		"https://www.flipkart.com/",
		"https://www.amazon.com/",
		"https://www.google.com",
		"https://www.flipkart.com/",
		"https://www.amazon.com/",
	}

	m := New(httpGetBody)
	for _, url := range urllist {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s %s %d bytes\n", url, time.Since(start), len(value.([]byte)))
	}
}
