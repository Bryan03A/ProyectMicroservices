package services

import (
	"context"
	"encoding/json"
	"log"

	"search-and-filters-service/config"
	"search-and-filters-service/models"

	"github.com/olivere/elastic/v7"
)

var esClient *elastic.Client

func init() {
	esConfig := config.GetElasticsearchConfig()
	client, err := elastic.NewClient(elastic.SetURL(esConfig.URL))
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	esClient = client
}

func Search(query string) ([]models.Document, error) {
	searchResult, err := esClient.Search().
		Index("documents").
		Query(elastic.NewQueryStringQuery(query)).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var docs []models.Document
	for _, hit := range searchResult.Hits.Hits {
		var doc models.Document
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func Index(doc models.Document) error {
	_, err := esClient.Index().
		Index("documents").
		BodyJson(doc).
		Do(context.Background())
	return err
}
