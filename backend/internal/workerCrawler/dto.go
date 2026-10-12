package workerCrawler

type crawlRequest struct{
	URL string `json:"url"`
}

type crawlResponse struct{
	URLs map[string][]string `json:"urls"`
}