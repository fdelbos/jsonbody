package main

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
)

type (
	Person struct {
		Name string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	PingHandler struct{}
)

var validate = validator.New()

const jsonBodyKey = "JSONBody"

func ExtractJSONBody(body interface{}) func(http.Handler) http.Handler {
	bodyType := reflect.TypeOf(body)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			instance := reflect.New(bodyType).Interface()

			err := json.NewDecoder(r.Body).Decode(instance)
			if err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if err := validate.Struct(instance); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), jsonBodyKey, instance)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetJSONBody(r *http.Request) interface{} {
	return r.Context().Value(jsonBodyKey)
}

func WriteJSONBody(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Print(err)
	}
}

func (PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := GetJSONBody(r).(*Person)
	// do something with it...

	log.Print(spew.Sdump(body))
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", ExtractJSONBody(Person{})(PingHandler{})).Methods("POST")
	return r
}

func main() {
	srv := &http.Server{
		Handler: Router(),
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
