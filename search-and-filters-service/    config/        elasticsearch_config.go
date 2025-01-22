// config/elasticsearch_config.go
package config

import (
    "github.com/olivere/elastic/v7"
    "log"
)

var (
    ElasticsearchHost = "your-elasticsearch-endpoint"
    ElasticsearchPort = 9200
    EsClient           *elastic.Client
)

func init() {
    client, err := elastic.NewClient(elastic.SetURL(elastic.SetURL(Scheme("http"), ElasticsearchHost, ElasticsearchPort)))
    if err != nil {
        log.Fatalf("Error creating the client: %s", err)
    }
    EsClient = client
}
