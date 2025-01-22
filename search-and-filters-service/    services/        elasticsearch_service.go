// services/elasticsearch_service.go
package services

import (
	"context"
	"encoding/json"
	"search-and-filters-service/config"

	"github.com/olivere/elastic/v7"
)

func IndexModel(modelID string, modelData map[string]interface{}) error {
	_, err := config.EsClient.Index().
		Index("modelos_3d").
		Id(modelID).
		BodyJson(modelData).
		Do(context.Background())
	return err
}

func SearchModels(query string) ([]map[string]interface{}, error) {
	searchResult, err := config.EsClient.Search().
		Index("modelos_3d").
		Query(elastic.NewMatchQuery("descripcion", query)).
		From(0).Size(10).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		var result map[string]interface{}
		err := json.Unmarshal(hit.Source, &result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
