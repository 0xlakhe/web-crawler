package workerCrawler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct{
	service *Service
}

func NewHandler(service *Service) *Handler{
	return &Handler{service: service}
}

func (h *Handler) CrawlHandler(w http.ResponseWriter, r *http.Request){
	
	var req crawlRequest
	if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{

		msg:=fmt.Sprintf("got error decoding: %v",err)
		_=json.NewEncoder(w).Encode(map[string]string{"error":msg})
		return
	}
	//request crawl
	resp,err:=h.service.Crawl(req)
	
	if err!=nil{
		msg:=fmt.Sprintf("got error crawling: %v",err)
		_=json.NewEncoder(w).Encode(msg)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

