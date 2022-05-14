package service

import (
    "mime/multipart"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)

func SearchPostsByUser(user string) ([]model.Post, error) {
	//array of Post, which is defined in model package.
    query := elastic.NewTermQuery("user", user)
	//first is the property time, second is the input user information
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	//backend is package, middle is global parameter, last is function
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
    query := elastic.NewMatchQuery("message", keywords)
    query.Operator("AND")
	//this and means you want all keywords must be included in the result
	//input is a string, not []string. Because multi keywords can be expressed as "xxx+yyy+zzz".
    if keywords == "" {
        query.ZeroTermsQuery("all")
    }
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
    var ptype model.Post
    var posts []model.Post

    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		//
        p := item.(model.Post)
        posts = append(posts, p)
    }
    return posts
}

func SavePost(post *model.Post, file multipart.File) error {
    medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
    if err != nil {
        return err
    }
    post.Url = medialink

    return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
    //what if this returns err so that es and gcs are not consistent?
    //if we do rollback, service is slowed down.
    //a good way: create a offline service that only runs once in a while to handle this.
}

func DeletePost(id string, user string) error {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("id", id))
    query.Must(elastic.NewTermQuery("user", user))

    return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
