package middleware

import (
	"net/http"
	"github.com/mileusna/useragent"
	"context"
)

type ctxKey struct{}

var osKey ctxKey

func SetOSInfo(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), osKey, ua.OS)
		r=r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetOSInfo(ctx context.Context) string {
	if v, ok := ctx.Value(osKey).(string); ok {
		return v
	}
	return ""
}

