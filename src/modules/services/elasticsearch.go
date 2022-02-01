package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
)

var ElasticSearchClient *elasticsearch.Client = nil

func InitializeElasticSearch(elasticHost string) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticHost},
	})
	if err != nil {
		fmt.Println("Failed to connect with elasticsearch", err)
		return
	}

	ElasticSearchClient = es
}

func IndexMessageInElasticSearch(index string, message []byte) {
	es := ElasticSearchClient

	if es == nil {
		log.Fatalf("ElasticSearch Client has not initialized")
		return
	}

	// Index the message in elastic search
	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(message),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return
	}

	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
		return
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return
	}

	// Print the response status and indexed document version.
	log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
}
