package workerCrawler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"sync"

	"golang.org/x/net/html"
)

type crawlResult struct {
	baseURL string
	urls    []string
	err     error
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func workers(receiveURL <-chan string, sendURL chan<- crawlResult) {

	for url := range receiveURL {
		links, err := downloadLink(url)
		if err != nil {
			sendURL <- crawlResult{url, nil, err}
		}
		sendURL <- crawlResult{url, links, nil}
	}
}

func (s *Service) Crawl(ctx context.Context,req crawlRequest) (*crawlResponse, error) {
	now:=time.Now()
	receiveURL := make(chan string, 50)
	sendURL := make(chan crawlResult)
	receiveURL <- req.URL
	allURLs := map[string][]string{}

	allLinks := map[string]int{}
	allLinks[req.URL] = 1
	var wg sync.WaitGroup

	maxWorkers := 5
	for range maxWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workers(receiveURL, sendURL)
		}()
	}

	go func() {
		wg.Wait()
		close(sendURL)
	}()
	crawlLimit := 20
	check := false
	for urls := range sendURL {
		if urls.err != nil {
			log.Print(urls.err)
			continue
		}
		toAdd := []string{}
		for _, url := range urls.urls {
			if _, ok := allLinks[url]; !ok {
				allLinks[url] = 1
				toAdd = append(toAdd, url)
				crawlLimit -= 1
				if check == true {
					continue
				}
				if crawlLimit < 1 {
					close(receiveURL)
					check = true
					continue
				}

				receiveURL <- url
			}
		}
		allURLs[urls.baseURL] = toAdd
	}
	total:=time.Since(now)
	timeTaken:=fmt.Sprintf("%.2f secs",total.Seconds())
	return &crawlResponse{URLs: allURLs,TimeTaken: timeTaken}, nil
}

func downloadLink(baseURL string) ([]string, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	token := html.NewTokenizer(resp.Body)
	result := []string{}

	for {
		tokenType := token.Next()
		switch tokenType {
		case html.ErrorToken:
			if token.Err() == io.EOF {
				fmt.Println("reached end of file")
				return result, nil
			}
			return nil, fmt.Errorf("error: %v", token.Token().Data)
		case html.StartTagToken:
			t := token.Token()
			if t.Data == "a" {

				for _, attribute := range t.Attr {
					if attribute.Key == "href" {
						relativeURL, err := url.Parse(attribute.Val)

						if err != nil {
							continue
						}
						base, err := url.Parse(baseURL)
						if err != nil {
							continue
						}
						toAdd := base.ResolveReference(relativeURL)
						result = append(result, toAdd.String())
					}
				}
			}
		}
	}
}
