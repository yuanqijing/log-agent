package client

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/yuanqijing/log-agent/apis"
	"strings"
)

var ESClient = NewElasticsearch()

type Elasticsearch struct {
	client *elasticsearch.Client
}

func NewElasticsearch() *Elasticsearch {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://es-elasticsearch-master:9200"},
		Username:  "elastic",
		Password:  "Hc@Cloud01",
	})

	if err != nil {
		panic(err)
	}

	return &Elasticsearch{
		client: client,
	}
}

// ip : es-elasticsearch-master
// port : 9200

func (e *Elasticsearch) GetClient() *elasticsearch.Client {
	return e.client
}

func (e *Elasticsearch) Index(index string, data apis.Log) (*esapi.Response, error) {

	dataByte, err := data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return e.client.Index(
		index,
		strings.NewReader(string(dataByte)),
	)
}
