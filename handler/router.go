package handler
//router is like dispatch servlet
import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	//middleware is put between router and handler
	//multi middlewares can be put to judge different things
    jwt "github.com/form3tech-oss/jwt-go"
	//generate secured token

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {

    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(mySigningKey), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })

    router := mux.NewRouter()

    router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
    router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")

    router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

    return router
}
