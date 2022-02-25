package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"

	"github.com/Jimeux/elastic-go/address"
)

var es *elasticsearch.Client

func init() {
	es, _ = elasticsearch.NewClient(elasticsearch.Config{
		// https://www.elastic.co/blog/the-go-client-for-elasticsearch-configuration-and-customization
		Addresses:         []string{"http://localhost:9210"},
		EnableDebugLogger: true,
	})
	/*res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)*/
}

func main() {
	ctx := context.Background()
	createIndex(ctx)
	// deleteIndex(ctx)
	// index(ctx)
	// update(ctx)
	// search(ctx)
}

func createIndex(ctx context.Context) {
	// Spec: https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-put-template.html
	// Kibana: https://www.elastic.co/guide/en/elasticsearch/reference/master/index-mgmt.html

	// create index template
	res, err := es.Indices.PutIndexTemplate(
		address.TemplateName,
		strings.NewReader(address.Template),
		es.Indices.PutIndexTemplate.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("cannot create template: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create template: %s", res)
	}

	// create index
	res, err = es.Indices.Create(address.IndexName, es.Indices.Create.WithContext(ctx))
	if err != nil {
		log.Fatalf("cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot create index: %s", res)
	}

	// create alias
	res, err = es.Indices.PutAlias([]string{address.IndexName}, address.Alias)
	if err != nil {
		log.Fatalf("cannot create alias: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot create alias: %s", err)
	}
}

func deleteIndex(ctx context.Context) {
	res, err := es.Indices.Delete([]string{address.IndexName}, es.Indices.Delete.WithContext(ctx))
	if err != nil {
		log.Fatalf("cannot delete index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot delete index: %s", res)
	}
	res, err = es.Indices.DeleteIndexTemplate(address.TemplateName, es.Indices.DeleteIndexTemplate.WithContext(ctx))
	if err != nil {
		log.Fatalf("cannot delete template: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot delete template: %s", res)
	}
}

func index(ctx context.Context) {
	// https://www.elastic.co/blog/the-go-client-for-elasticsearch-working-with-data

	doc := address.Address{
		ID:    1,
		Line1: "ユースフル幡ヶ谷",
		Line2: "渋谷区",
	}

	res, err := es.Index(address.Alias, esutil.NewJSONReader(&doc),
		es.Index.WithContext(ctx),
		es.Index.WithRequireAlias(true),
		es.Index.WithDocumentID(strconv.FormatInt(doc.ID, 10)),
	)
	if err != nil {
		log.Fatalf("cannot insert doc: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot insert doc: %s", res)
	}
}

func update(ctx context.Context) {
	// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-update.html

	body := map[string]any{
		"doc": map[string]any{
			"line_1": "ラ・メイゾン",
		},
	}

	res, err := es.Update(address.Alias, "1", esutil.NewJSONReader(&body),
		es.Update.WithContext(ctx),
		es.Update.WithRequireAlias(true),
	)
	if err != nil {
		log.Fatalf("cannot insert doc: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot insert doc: %s", res)
	}
}

func search(ctx context.Context) {
	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"line_2": "渋谷区",
			},
		},
	}

	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(address.Alias),
		es.Search.WithBody(esutil.NewJSONReader(&query)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("cannot search: %s", err)
	}
	if res.IsError() {
		log.Fatalf("cannot search: %s", res)
	}

	var e map[string]any
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", e)
}
