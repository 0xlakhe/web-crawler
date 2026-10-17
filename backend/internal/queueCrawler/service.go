//Implemented crawler using queue

package queueCrawler

import (
	"errors"
	"fmt"
	"time"

	downloadlink "github.com/0xlakhe/web-scrapper-go.git/internal/downloadLink"

)

var errorQueueEmpty = errors.New("queue is empty")

type queue[T any] struct {
	Urls []T
}

func (q *queue[T]) enqueue(toQueue T) {
	q.Urls = append(q.Urls, toQueue)
}

func (q *queue[T]) dequeue() (T, error) {
	var zero T
	if len(q.Urls) == 0 {
		return zero, errorQueueEmpty
	}
	res := q.Urls[0]
	q.Urls[0] = zero
	q.Urls = q.Urls[1:]
	return res, nil
}

type Service struct{

}

func NewService()*Service{
	return &Service{}
}

func (s *Service)Crawl(req crawlRequest)(*crawlResponse,error) {
	now:=time.Now()
	baseUrl := req.URL
	htmlUrls := map[string][]string{}
	queue := queue[string]{}
	queue.enqueue(baseUrl)
	allUrls := map[string]int{}
	allUrls[baseUrl] = 1

	//crawl limit
	maxPages := 20
	currentPage := 1
	for currentPage < maxPages {
		baseURL, err := queue.dequeue()
		if err != nil {
			if errors.Is(err, errorQueueEmpty) {
				return nil,err
			} else {
				
				return nil,err
			}
		}
		urls, err := downloadlink.Download(baseURL)

		allLinks := []string{}
		//adding links to map
		for _, link := range urls {

			//adding link to queue
			if _, ok := allUrls[link]; !ok {
				allLinks = append(allLinks, link)
				queue.enqueue(link)
			}
			allUrls[link] = 1
		}

		if err != nil {
			return nil,err
		}

		//saving urls
		htmlUrls[baseURL] = allLinks
		currentPage += 1
	}
	total:=time.Since(now)
	timeTaken:=fmt.Sprintf("%.2f secs",total.Seconds())
	return &crawlResponse{URLs: htmlUrls, TimeTaken: timeTaken},nil

}
