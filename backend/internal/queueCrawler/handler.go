package queueCrawler

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

func (h *Handler)CrawlHandler(w http.ResponseWriter, r *http.Request){
	var req crawlRequest
	err:=json.NewDecoder(r.Body).Decode(&req)

	w.Header().Set("Content-Type","application/json")
	if err!=nil{
		e:=fmt.Sprintf("error decoging: %v\n",err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e))
		return
	}

	resp,err:=h.service.Crawl(req)
	if err!=nil{
		e:=fmt.Sprintf("error crawling: %v\n",err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e))
		return
	}
	_=json.NewEncoder(w).Encode(resp)
}