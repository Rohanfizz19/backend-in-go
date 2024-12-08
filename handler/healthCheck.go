package handler

import (
	"backend/httpserver"
	"net/http"
)

func NewHealthChecker(router httpserver.Router) {
	router.AddRoute(httpserver.RouteConfig{
		Path: "/rest/health",
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
			resp.Write([]byte("success"))
		}),
		Methods:    []string{http.MethodGet},
		Instrument: true,
	})
}
