package httpserver

import (
	"backend/config"
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	ContextKeyReqID                         string = "requestId"
	ContextKeyCustomerID                    string = "customerId"
	ContextKeyDeviceID                      string = "deviceId"
	ContextKeySource                        string = "source"
	HTTPHeaderNameRequestID                 string = "X-Request-ID"
	HTTPHeaderNameAccessControlAllowOrigin  string = "Access-Control-Allow-Origin"
	HTTPHeaderNameAccessControlAllowMethods string = "Access-Control-Allow-Methods"
	HTTPHeaderNameAccessControlAllowHeaders string = "Access-Control-Allow-Headers"
)

type Middleware func(next http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) < 1 {
		return h
	}

	wrapper := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapper = middlewares[i](wrapper)
	}

	return wrapper
}

func CorsMiddleware(cors *config.Cors) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header()[HTTPHeaderNameAccessControlAllowOrigin] = cors.Origins
			w.Header()[HTTPHeaderNameAccessControlAllowMethods] = cors.Methods
			w.Header()[HTTPHeaderNameAccessControlAllowHeaders] = cors.Headers

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := attachReqID(r.Context())
		r = r.WithContext(ctx)
		w.Header().Set(HTTPHeaderNameRequestID, GetReqID(ctx))

		next.ServeHTTP(w, r)
	})
}
func attachReqID(ctx context.Context) context.Context {
	reqID := genReqId()
	return AttachToCtx(ctx, ContextKeyReqID, reqID)
}

func GetReqID(ctx context.Context) string {
	reqID := ctx.Value(ContextKeyReqID)
	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}
func genReqId() string {
	rqId, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return rqId.String()
}

func AttachToCtx(ctx context.Context, key string, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
