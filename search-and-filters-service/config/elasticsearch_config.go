package config

import "log"

type ElasticsearchConfig struct {
	URL string `mapstructure:"url"`
}

func GetElasticsearchConfig() ElasticsearchConfig {
	var elasticsearchConfig ElasticsearchConfig
	if err := Config.UnmarshalKey("elasticsearch", &elasticsearchConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return elasticsearchConfig
}
