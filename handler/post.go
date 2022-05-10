package handler

import (
	"encoding/json" //similar to jackson
	"fmt"           //format, like print
	"net/http"

	"around/model"
	"around/service"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse from body of request to get a json object.
	// the pointer, or the input object will be copied and the original object can't be modified
	// why to change resuest, but not response?
	// - request pointer is space saving. you are not going to change the request, just to save space
	// responseWriter is an interface, not a class. There is no pointer to interfaces
	fmt.Println("Received one post request")
	decoder := json.NewDecoder(r.Body)
	// here you must know request body is json
	var p model.Post
	if err := decoder.Decode(&p); err != nil {
		//json is decoded to struct p
		panic(err)
		//panic = throw a runtime exception
		//everything crash and then restart.
		//normally you don't do this.
	}

	fmt.Fprintf(w, "Post received: %s\n", p.Message)
	//print these to w
	//writer is like a buffer. print to writer first, and writer write to response.
	//writer will judge if the string is too long for response or something else.
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords") //query in url is the part behind ?

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		//make it a http response. third is the error name, it is 500.
		//2++ ok
		//3++ can't find here, but know how to find it
		//4++ error from client
		//5++ error from server
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	js, err := json.Marshal(posts) //work with all go structs
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
