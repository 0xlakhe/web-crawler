package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/0xlakhe/web-scrapper-go.git/internal/queueCrawler"
	"github.com/0xlakhe/web-scrapper-go.git/internal/workerCrawler"
)

func main() {
	


	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	})

	workerModule := workerCrawler.NewModule()
	mux.HandleFunc("POST /worker-crawl", workerModule.Handler.CrawlHandler)

	queueModule:=queueCrawler.NewModule()
	mux.HandleFunc("POST /queue-crawl",queueModule.Handler.CrawlHandler)
	check := make(chan string)

	go func() {
		http.ListenAndServe(":8080", mux)

	}()
	fmt.Println("running in port 8080")
	<-check
}
