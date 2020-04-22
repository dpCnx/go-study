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
	//quary()
}

/*
	初始es
*/
func inites() {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://10.25.26.195:9200",
		},
		Transport: &http.Transport{

			// MaxIdleConnsPerHost，如果非零，则控制最大空闲 (keep-alive)保持每个主机的连接。
			//如果为零, 使用DefaultMaxIdleConnsPerHost
			MaxIdleConnsPerHost: 10,
			// 如果非零，指定的数量 等待服务器响应报头的时间 编写请求(包括请求正文，如果有的话)。
			//这时间不包括读取响应体的时间。
			ResponseHeaderTimeout: 30 * time.Second,
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

	/*res, err := esclient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)*/
}

/*
	获取节点
*/
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

/*
	esapi.CreateRequest
*/
func insert() {

	/*	student := &Student{
			Name:    "dp",
			Age:     18,
			Address: "重庆",
		}
		jsonBody, _ := json.Marshal(student)*/

	body := map[string]interface{}{
		"name": "dp",
		"age":  18,
	}
	jsonBody, _ := json.Marshal(body)

	req := esapi.CreateRequest{ // 如果是esapi.IndexRequest则是插入/替换
		Index:        "demo",
		DocumentType: "test1",
		DocumentID:   "1",
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

/*
	esclient.Create
*/
func insert2() {

	student := &Student{
		Name:    "dp",
		Age:     18,
		Address: "重庆",
	}
	jsonBody, _ := json.Marshal(student)

	res, err := esclient.Create("demo", "4", bytes.NewReader(jsonBody), esclient.Create.WithDocumentType("test1"))
	if err != nil {
		log.Printf("create err:%v", err)
		return
	}

	log.Println(res)
}

/*
	esapi.BulkRequest
*/
func insert3() {

	var data bytes.Buffer

	createLine := map[string]interface{}{
		"create": map[string]interface{}{
			"_id": "5",
		},
	}
	jsonStr, _ := json.Marshal(createLine)
	data.Write(jsonStr)
	data.WriteByte('\n')

	student1 := &Student{
		Name:    "dp",
		Age:     19,
		Address: "重庆",
	}

	jsonBody, _ := json.Marshal(&student1)
	data.Write(jsonBody)
	data.WriteByte('\n')

	createLine = map[string]interface{}{
		"create": map[string]interface{}{
			"_id": "6",
		},
	}
	jsonStr, _ = json.Marshal(createLine)
	data.Write(jsonStr)
	data.WriteByte('\n')

	student2 := &Student{
		Name:    "dp",
		Age:     19,
		Address: "重庆",
	}

	jsonBody, _ = json.Marshal(&student2)
	data.Write(jsonBody)
	data.WriteByte('\n')

	req := esapi.BulkRequest{
		Index:        "demo",
		DocumentType: "test1",
		Body:         &data,
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("insert err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())

}

/*
	esclient.Bulk
*/
func insert4() {
	/*	res, err := esclient.Bulk(
				strings.NewReader(`
		{ "index" : { "_index" : "demo","_type":"test1", "_id" : "7" } }
		{ "name" : "cnx" ,"age":1}
		{ "index" : { "_index" : "demo","_type":"test1", "_id" : "8" } }
		{ "name" : "cnx2" ,"age":2}
		{ "index" : { "_index" : "demo","_type":"test1", "_id" : "9" } }
		{ "name" : "cnx3" ,"age":3}
		`), )

			//bytes.NewReader()
			if err != nil {
				log.Println("bulk err:", err)
				return
			}

			log.Println(res)*/

	res, err := esclient.Bulk(
		strings.NewReader(`
{ "index" : { "_id" : "10" } }
{ "name" : "cnx" ,"age":12}
`), esclient.Bulk.WithIndex("demo"), esclient.Bulk.WithDocumentType("test1"))

	//bytes.NewReader()
	if err != nil {
		log.Println("bulk err:", err)
		return
	}

	log.Println(res)
}

/*
	esapi.UpdateRequest
*/
func update() {
	body := map[string]interface{}{
		"doc": map[string]interface{}{
			"name": "cnx",
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.UpdateRequest{
		Index:        "demo",
		DocumentType: "test1",
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
	esapi.UpdateRequest
*/
func update2() {

	body := map[string]interface{}{
		"doc": map[string]interface{}{
			"name": "d",
		},
	}
	jsonBody, _ := json.Marshal(body)

	res, err := esclient.Update("demo", "1", bytes.NewReader(jsonBody), esclient.Update.WithDocumentType("test1"))
	if err != nil {
		log.Println("update err:", err)
		return
	}
	log.Println(res)
}

/*
	esapi.DeleteRequest
*/
func delete() {
	req := esapi.DeleteRequest{
		Index:        "demo",
		DocumentType: "test1",
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

/*
	esapi.DeleteByQueryRequest
*/
func delete2() {
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "d",
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req := esapi.DeleteByQueryRequest{
		Index:        []string{"demo"},
		DocumentType: []string{"test1"},
		Body:         bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("delete err :", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

/*
	esapi.IndicesDeleteRequest
*/
func deleteIndex() {

	req := esapi.IndicesDeleteRequest{
		Index: []string{"test1"},
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("delete err :", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())

}

/*
	esapi.IndicesCreateRequest{}
*/
func creadIndex() {

	//elasticsearch7默认不在支持指定索引类型，默认索引类型是_doc
	//elasticsearch6 可以正常执行
	/*body := map[string]interface{}{
		"mappings": map[string]interface{}{
			"logs": map[string]interface{}{
				"properties": map[string]interface{}{
					"title": map[string]string{
						"type": "text",
					},
					"author": map[string]string{
						"type": "text",
					},
					"titleScore": map[string]string{
						"type": "double",
					},
				},
			},
		},
	}*/

	body := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]string{
					"type": "text",
				},
				"author": map[string]string{
					"type": "text",
				},
				"titleScore": map[string]string{
					"type": "double",
				},
			},
		},
	}

	b, _ := json.Marshal(body)

	log.Println(string(b))

	req := esapi.IndicesCreateRequest{
		Index: "demo3",
		Body:  bytes.NewReader(b),
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("create err :", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

/*
	esclient.Indices.Create -->未成功
*/
func createIndex2() {
	res, err := esclient.Indices.Create(
		"demo1",
		esclient.Indices.Create.WithBody(strings.NewReader(`{
    "settings":{
        "number_of_shards":3,
        "number_of_replicas":2
    },
    "mappings":{
        "properties":{
            "name":{
                "type":"text"
            }
        }
    }
}`)),
	)
	fmt.Println(res, err)
}

func quaryone() {
	req := esapi.GetRequest{
		Index:        "demo",
		DocumentType: "test1",
		DocumentID:   "2",
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())
}

func quarym() {
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
		Index:        "demo",
		DocumentType: "test1",
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

/*
	不加条件就是查询所有
*/
func quary() {

	/*body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]string{
				"about": "travel",
			},
		},
	}*/

	/*body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": map[string]string{
						"about": "travel",
					},
				},
				"must_not": map[string]interface{}{
					"match": map[string]string{
						"sex": "boy",
					},
				},
			},
		},
	}*/

	/*body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"terms": map[string]interface{}{
						"about": []string{"travel", "history"},
					},
				},
			},
		},
	}*/

	/*body := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"age": map[string]interface{}{
					"gt":  20,
					"lte": 25,
				},
			},
		},
	}*/

	/*body := map[string]interface{}{
		"query": map[string]interface{}{
			"exists": map[string]interface{}{
				"field": "age",
			},
		},
	}*/

	/*body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"about": map[string]interface{}{
								"value": "travel",
							},
						},
					},
					{
						"term": map[string]interface{}{
							"name": "daqiao",
						},
					},
					{
						"range": map[string]interface{}{
							"age": map[string]interface{}{
								"gte": 20,
								"lte": 30,
							},
						},
					},
				},
			},
		},
	}*/

	/*	body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": map[string]interface{}{
						"sex": "girl",
					},
				},
				"filter": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"age": 20,
						},
					},
				},
			},
		},
	}*/

	body := map[string]interface{}{
		"from": 0,
		"size": 1,
	}

	b, _ := json.Marshal(body)

	log.Println(string(b))

	req := esapi.SearchRequest{
		Index:        []string{"demo"},
		DocumentType: []string{"student"},
		Body:         bytes.NewReader(b),
	}

	res, err := req.Do(context.Background(), esclient)
	if err != nil {
		log.Println("get err:", err)
		return
	}
	defer res.Body.Close()
	log.Println(res.String())

}

func search() {
	res, err := esclient.Search(
		esclient.Search.WithIndex("demo"),
		//esclient.Search.WithDocumentType("test1"),
		//esclient.Search.WithSort("age:desc"),
		esclient.Search.WithQuery("name:c*"), // *：匹配任意多个字符  ？：仅匹配一个字符
		//esclient.Search.WithScroll(3),
		//esclient.Search.WithSize(1),

		esclient.Search.WithPretty(), //格式化
	)

	fmt.Println(res, err)
}

//https://my.oschina.net/u/3100849/blog/1839022  类型

//dynamic 新增字段情况，Dynamic 设置为 true，带有新字段的文档写入，Mapping 会更新。
// 						Dynamic 设置为 false，Mapping 不被更新，新增字段不会被索引。
// 						Dynamic 设置为 Strict，带有新字段的文档写入会直接报错

//esapi.ScrollRequest{}
