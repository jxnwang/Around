package handler

import (
	"encoding/json" //similar to jackson
	"fmt"           //format, like print
	"net/http"
	"path/filepath"

	"around/model"
	"around/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/pborman/uuid"
	"github.com/gorilla/mux" 
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")

	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	p := model.Post{
		Id:   uuid.New(), //creaqte a unique id
		User: username.(string),
		//user identity should be from the token, not given by http request as message.
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	//second return is metadata
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename) //.mp4, .jpg ...
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	fmt.Println("Post is saved successfully.")

	// Parse from body of request to get a json object.
	// the pointer, or the input object will be copied and the original object can't be modified
	// why to change resuest, but not response?
	// - request pointer is space saving. you are not going to change the request, just to save space
	// responseWriter is an interface, not a class. There is no pointer to interfaces
	//	fmt.Println("Received one post request")
	//	decoder := json.NewDecoder(r.Body)
	// here you must know request body is json
	//	var p model.Post
	//	if err := decoder.Decode(&p); err != nil {
	//json is decoded to struct p
	//		panic(err)
	//panic = throw a runtime exception
	//everything crash and then restart.
	//normally you don't do this.
	//	}

	//	fmt.Fprintf(w, "Post received: %s\n", p.Message)
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete post from backend %v\n", err)
        return
    }
    fmt.Println("Post is deleted successfully")
}
