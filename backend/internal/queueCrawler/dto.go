package queueCrawler

type crawlRequest struct{
	URL string `json:"url"`
}

type crawlResponse struct{
	URLs map[string][]string `json:"urls"`
	TimeTaken string `json:"time_taken"`
}