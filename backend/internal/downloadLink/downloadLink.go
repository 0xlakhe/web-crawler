package downloadlink

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)


func Download(baseURL string) ([]string, error) {
	
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
