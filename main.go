package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

func main() {
	http.HandleFunc("/infonpm", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.ListenAndServe(":8080", nil)
}