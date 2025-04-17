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