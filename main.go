package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	Person struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	PingHandler struct{}

	contextKey int
)

var validate = validator.New()

const jsonBodyKey contextKey = iota

func JSONBody(body interface{}) func(http.Handler) http.Handler {
	bodyType := reflect.TypeOf(body)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			switch r.Header.Get("Content-Type") {
			case "application/json", "application/json; charset=utf-8":
			default:
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			instance := reflect.New(bodyType).Interface()
			err := json.NewDecoder(r.Body).Decode(instance)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if err := validate.Struct(instance); err != nil {
				// should display a formatted error instead...
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

func (p PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := GetJSONBody(r).(*Person)
	WriteJSONBody(w, body)
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", JSONBody(Person{})(PingHandler{})).Methods(http.MethodPost)
	return r
}

func main() {
	srv := &http.Server{
		Handler: Router(),
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
