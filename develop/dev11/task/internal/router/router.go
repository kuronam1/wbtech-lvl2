package router

import "net/http"

type AppRouter struct {
	router *http.ServeMux
}

func NewRouter() http.ServeMux {
	router := http.NewServeMux()

}
