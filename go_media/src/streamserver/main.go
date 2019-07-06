package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHanlders() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)

	return router
}

func main() {
	router := RegisterHanlders()
	http.ListenAndServe(":9000", router)
}
