package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	esclient *elasticsearch.Client
	err      error
)

func main() {
	inites()
	//getnodes()
	//creadindex()
	//insertdocument()
	//getone()
	//getall()
	//getmrequest()
	//getbysearch()
	//update()
	//updateByQuery()
	//delete()
	//deleteByQuery()
	//bulk()
	//bulkrequest()
}

func bulkrequest() {
	req := esapi.BulkRequest{
		Index:        "demo1",
		DocumentType: "demo1_test",
		Body: strings.NewReader(`
{ "index" : {"_id" : "7" }}
{ "name" : "dp6" ,"age":3}
`),
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("bulkrequest err :", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

/*
	批量导入
*/
func bulk() {
	res, err := esclient.Bulk(
		strings.NewReader(`
{ "index" : { "_index" : "demo1","_type":"demo1_test", "_id" : "4" } }
{ "name" : "cnx" ,"age":1}
{ "index" : { "_index" : "demo1","_type":"demo1_test", "_id" : "5" } }
{ "name" : "cnx2" ,"age":2}
{ "index" : { "_index" : "demo1","_type":"demo1_test", "_id" : "6" } }
{ "name" : "cnx3" ,"age":3}
`),
	)

	if err != nil {
		log.Println("bulk err:", err)
		return
	}

	log.Println(res)
}

func deleteByQuery() {
	body := map[string]interface{}{
		"query": map[string]interface{}{
			//"match_all": map[string]interface{}{}, //匹配所有
			"match": map[string]interface{}{
				"name": "dp2",
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.DeleteByQueryRequest{
		Index: []string{"demo1"},
		Body:  bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("delete err :", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func delete() {
	req := esapi.DeleteRequest{
		Index:        "demo1",
		DocumentType: "demo1_test",
		DocumentID:   "1",
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())

}

func update() {
	body := map[string]interface{}{
		"doc": map[string]interface{}{
			"name": "cnx",
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.UpdateRequest{
		Index:        "demo1",
		DocumentType: "demo1_test",
		DocumentID:   "1",
		Body:         bytes.NewReader(jsonBody),
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

/*
	暂时不能使用
*/
func updateByQuery() {
	body := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.likes++",
		},
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"age": 16,
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.UpdateByQueryRequest{
		Index: []string{"demo1"},
		Body:  bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("updata err:", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func getbysearch() {
	res, err := esclient.Search(
		esclient.Search.WithIndex("demo1/demo1_test"),
		//esclient.Search.WithQuery("_id:3"), //精确查询
		//esclient.Search.WithSize(1),
		esclient.Search.WithPretty(),
	)
	fmt.Println(res, err)
}

func getall() {
	req := esapi.SearchRequest{
		Index:        []string{"demo1"},
		DocumentType: []string{"demo1_test"},
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())
}

func getmrequest() {
	body := map[string]interface{}{
		"docs": []map[string]interface{}{
			{
				"_id": "3",
			},
			{
				"_id": "4",
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.MgetRequest{
		Index:        "demo1",
		DocumentType: "demo1_test",
		Body:         bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("mgetrequest err :", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func getone() {
	req := esapi.GetRequest{
		Index:        "demo1",
		DocumentType: "demo1_test",
		DocumentID:   "1",
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())
}

func insertdocument() {
	body := map[string]interface{}{
		"name": "dp3",
		"age":  18,
	}
	jsonBody, _ := json.Marshal(body)

	req := esapi.CreateRequest{ // 如果是esapi.IndexRequest则是插入/替换
		Index:        "demo1",
		DocumentType: "demo1_test",
		DocumentID:   "3",
		Body:         bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("insert err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())
}

func creadindex() {
	res, err := esclient.Index(
		"test",                                                      // Index name
		strings.NewReader(`
	{
		"title" : "Test"
	}
	`), // Document body
		esclient.Index.WithDocumentID("1"),                          // Document ID
		esclient.Index.WithRefresh("true"),                          // Refresh
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
}

func getnodes() {
	nodesRequest := esapi.CatNodesRequest{}
	respose, err := nodesRequest.Do(context.Background(), esclient)
	if err != nil {
		log.Println("nodes err:", err)
		return
	}

	defer respose.Body.Close()

	log.Println(respose)
}

//初始化es
func inites() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.0.108:9200",
		},
		Transport: &http.Transport{

			// MaxIdleConnsPerHost，如果非零，则控制最大空闲 (keep-alive)保持每个主机的连接。
			//如果为零, 使用DefaultMaxIdleConnsPerHost
			MaxIdleConnsPerHost: 10,
			// 如果非零，指定的数量 等待服务器响应报头的时间 编写请求(包括请求正文，如果有的话)。
			//这时间不包括读取响应体的时间。
			ResponseHeaderTimeout: time.Second,
			//DialContext与往返电话同时运行。
			//发起拨号的往返呼叫可能最终使用
			//以前连接时所拨打的连接
			//在后面的拨号连接完成之前就空闲了。
			DialContext: (&net.Dialer{Timeout: time.Second}).DialContext,

			// TLSClientConfig指定要使用的TLS配置
			// tls.Client。如果为空，则使用默认配置。
			//如果非nil, HTTP/2支持可能默认不启用。
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}

	esclient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("init es err :", err)
		return
	}

	/*	res, err := esclient.Info()
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		log.Println(res)*/
}
