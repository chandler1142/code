package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r  *httprouter.Router
	cl *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.cl.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	m.r.ServeHTTP(w, r)
	defer m.cl.ReleaseConn()
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.cl = NewConnLimiter(cc)
	return m
}

func RegisterHanlders() *httprouter.Router {
	router := httprouter.New()

	router.GET("/testpage", testPageHandler)

	router.GET("/video/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)

	return router
}

func main() {
	router := RegisterHanlders()
	mw := NewMiddleWareHandler(router, 2)
	http.ListenAndServe(":9000", mw)
}
