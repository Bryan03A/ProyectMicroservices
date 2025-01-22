// utils/search_utils.go
package utils

import "github.com/olivere/elastic/v7"

func BuildQuery(query string, filters map[string]interface{}) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery("descripcion", query))
	if filters != nil {
		boolQuery.Filter(filters)
	}
	return boolQuery
}
