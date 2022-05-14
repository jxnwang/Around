package backend

import (
    "context"
    "fmt"

    "around/constants"
    "around/util"

    "github.com/olivere/elastic/v7"
)

var (
    ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
    client *elastic.Client
}

func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
    client, err := elastic.NewClient(
        elastic.SetURL(config.Address),
        elastic.SetBasicAuth(config.Username, config.Password))
    if err != nil {
        panic(err)
    }

    exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        //keyword: you can only use strict mapping
        //select * from post where id = "123"
        //text: select * ... contains "..."

        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
        //"index": false. this index is not a db, but to mark if a better search method id added to this coloum.
        //you can still search for something, but time will be O(N)
        _, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
        // this _ means there is a return value that you don't want
        if err != nil {
            panic(err)
        }
    }


    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")

    ESBackend = &ElasticsearchBackend{client: client}
    //first is the pointer dewfined in line 17. second is returned by function in line 21
    // : means initialize
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    // i is interface to support different saved items
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())


    return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    return err
}