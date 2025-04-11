package middlewares

import "net/http"

func CorsMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5137")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
        handler.ServeHTTP(w,r)
    }
}
