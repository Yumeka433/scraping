package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Masukan parameter query"})
		return
	}

	c := colly.NewCollector()
	url := fmt.Sprintf("https://registry.npmjs.org/%s", query)

	c.OnResponse(func(r *colly.Response) {
		var data map[string]interface{}
		err := json.Unmarshal(r.Body, &data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Error fetching data"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": 200, "result": data})
	})

	c.OnError(func(r *colly.Response, err error) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Error fetching data"})
	})

	c.Visit(url)
}

func main() {
	http.HandleFunc("/api/infonpm", handler)
	http.ListenAndServe(":8080", nil)
}
```
Namun, karena Vercel menggunakan serverless function, Anda perlu mengubah kode di atas menjadi serverless function. Berikut adalah contoh kode yang dapat digunakan:
```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	query := request.QueryStringParameters["query"]

	if query == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"status": false, "message": "Masukan parameter query"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	c := colly.NewCollector()
	url := fmt.Sprintf("https://registry.npmjs.org/%s", query)

	var data map[string]interface{}
	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &data)
		if err != nil {
			data = map[string]interface{}{"status": false, "message": "Error fetching data"}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		data = map[string]interface{}{"status": false, "message": "Error fetching data"}
	})

	c.Visit(url)

	jsonData, _ := json.Marshal(map[string]interface{}{"status": 200, "result": data})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
