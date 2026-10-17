# web-crawler

A small Go web crawler with a JSON HTTP API. Given a start URL, it fetches the page, extracts links, follows them (up to a page limit), and returns the discovered URLs.


## Getting started

All commands are run from the `backend/` directory, which holds the Go module
(`github.com/0xlakhe/web-scrapper-go.git`).

```bash
cd backend

go build ./...
go vet ./...
go run ./cmd/api
```

The API starts on `http://localhost:8080`.

## API

### `GET /hello`

Health/hello check. Returns `Hello`.

```bash
curl http://localhost:8080/hello
```

### `POST /queue-crawl`

Crawls starting from the given URL using a sequential breadth-first search over
an in-memory queue.

### `POST /worker-crawl`

Same job as `/queue-crawl`, but downloads pages using a pool of 5 concurrent
workers.

Both endpoints take the same JSON body and return the same shape:

```bash
curl -X POST http://localhost:8080/queue-crawl \
  -H "Content-Type: application/json" \
  -d '{"url": "https://web-scraping.dev/"}'
```

Request:

```json
{ "url": "https://example.com/" }
```

Response:

```json
{
  "urls": {
    "https://example.com/": ["https://example.com/about", "https://example.com/contact"],
    "https://example.com/about": ["https://example.com/"]
  },
  "time_taken": "1.24 secs"
}
```

`urls` maps each crawled page to the new links discovered on it. Each crawler is
limited to roughly 20 pages.

## Project layout

```
backend/
  cmd/api/                HTTP server and route wiring
  internal/queueCrawler/  queue-based crawler
  internal/workerCrawler/ worker-pool crawler
  internal/routines/      standalone scratch/playground, not part of the API
  text                    sample crawl output (JSON)
```

Each crawler package follows the same layout: `module.go` (wiring),
`handler.go` (HTTP), `service.go` (crawl logic), `dto.go` (request/response
types).
