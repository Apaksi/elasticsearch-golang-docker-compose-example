package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	time.Sleep(3 * time.Second)

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %s", err)
	}

	var dbUrl string
	if hostname == "app" {
		dbUrl = "https://elasticsearch:9200"
	} else {
		dbUrl = "https://localhost:9200"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{dbUrl},
		Username:  "elastic",
		Password:  "elastinen",
		Transport: &http.Transport{
			//MaxIdleConnsPerHost:   10,
			//ResponseHeaderTimeout: 5 * time.Second,
			//DialContext: (&net.Dialer{
			//	Timeout:   1 * time.Second,
			//	KeepAlive: 30 * time.Second,
			//}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				//	RootCAs:            rootCAs,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Insert documents
	index := "my_index"
	documents := []string{
		"This is the first document.",
		"This is the second document.",
		"This is the third document.",
	}

	for i, doc := range documents {
		putDoc := fmt.Sprintf(`{"text": "%s"}`, doc)
		req := esapi.IndexRequest{
			Index:      index,
			DocumentID: fmt.Sprintf("%d", i+1),
			Body:       io.NopCloser(strings.NewReader(putDoc)),
			Refresh:    "true",
		}
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error indexing document: %s", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("Error indexing document: %s", res.Status())
		} else {
			log.Printf("Document indexed successfully: %s", doc)
		}
	}

	// Fuzzy search
	query := `
	{
		"query": {
			"fuzzy": {
				"text": {
					"value": "documant",
					"fuzziness": "AUTO"
				}
			}
		}
	}
	`
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error searching documents: %s", res.Status())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		fmt.Println("Matched document:", source["text"])
	}
}
