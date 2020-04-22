package main

import (
	"bytes"
	"context"
	"encoding/json"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	esapi7 "github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
)

func main() {
	cfg := elasticsearch7.Config{
		Addresses: []string{
			"http://10.25.26.195:9200",
		},
	}
	esclient7, err := elasticsearch7.NewClient(cfg)
	if err != nil {
		log.Println("init err:", err)
		return
	}

	/*res, err := esclient7.Info()

	if err != nil {
		log.Println("info err:", err)
		return
	}

	log.Println(res)*/

	searchSql(esclient7)
}

func searchSql(es *elasticsearch7.Client) {

	body := map[string]interface{}{
		"query": "SELECT * FROM demo WHERE age = 18",
	}
	jsonBody, _ := json.Marshal(body)

	req := esapi7.SQLQueryRequest{
		Body:   bytes.NewReader(jsonBody),
		Format: "json",
		Pretty: true,
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())
}
