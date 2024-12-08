package handler

import (
	"backend/httpserver"
	"backend/metrics"
	"context"
	"fmt"
	"net/http"
	"time"
)

type LocationFeaturesAPIHandler struct {
	debugMode bool
}

type ContentGetUserListsAPIHandler struct {
	debugMode bool
}

func NewLocationAPIHandler(ctx context.Context, router httpserver.Router) {
	locationAPIHandler := &LocationFeaturesAPIHandler{}

	middlewares := []httpserver.Middleware{}

	wrapper := httpserver.ChainMiddleware(http.HandlerFunc(locationAPIHandler.ServeHTTP), middlewares...)
	debugWrapper := httpserver.ChainMiddleware(http.HandlerFunc(locationAPIHandler.ServeHTTPWithDebug), middlewares...)

	router.AddRoute(httpserver.RouteConfig{
		Path:       "/api/v1/location_based_features",
		Handler:    wrapper,
		Methods:    []string{http.MethodGet},
		Instrument: true,
	})
	router.AddRoute(httpserver.RouteConfig{
		Path:       "/api/v1/debug/location_based_features",
		Handler:    debugWrapper,
		Methods:    []string{http.MethodGet},
		Instrument: true,
	})
}

func (handler *LocationFeaturesAPIHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	defer metrics.InstrumentDataProviderLatencyHistogram("serveHTTP-LocationFeatureAPI", time.Now())
	// ctx := req.Context()

	// response := locationapi.Response{
	// 	Data:          data,
	// 	StatusCode:    0,
	// 	StatusMessage: "Done Successfully ",
	// 	DebugResponse: aggregatorResponse.DebugResponse,
	// 	DeviceID:      header.DeviceID,
	// 	TID:           header.TID,
	// 	SID:           header.SID,
	// }

	// handler.writeResponse(ctx, respWriter, &response)
}

func (handler *LocationFeaturesAPIHandler) ServeHTTPWithDebug(respWriter http.ResponseWriter, req *http.Request) {
	debugHandler := LocationFeaturesAPIHandler{
		debugMode: true,
	}
	debugHandler.ServeHTTP(respWriter, req)
}

func (handler *ContentGetUserListsAPIHandler) writeResponse(
	ctx context.Context,
	respWriter http.ResponseWriter,
	internalStatusCode int,
	statusMessage string,
) {

	_, err := respWriter.Write(([]byte)("asd"))
	if err != nil{
		fmt.Println(err)
	}
}
