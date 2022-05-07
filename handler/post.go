package handler

import (
    "encoding/json"//similar to jackson
    "fmt"//format, like print
    "net/http"
   
    "around/model"
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