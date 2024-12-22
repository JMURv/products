package middleware

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			span := opentracing.GlobalTracer().StartSpan(
				fmt.Sprintf("%s %s", r.Method, r.URL),
			)
			defer span.Finish()

			zap.L().Info("Request", zap.String("method", r.Method), zap.String("uri", r.RequestURI))
			next.ServeHTTP(w, r)
		},
	)
}
