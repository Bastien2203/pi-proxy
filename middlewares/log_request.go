package middlewares

import (
	"fmt"
	"net/http"
)

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("===============REQUEST=====================")
		fmt.Println("Host: ", r.Host)
		fmt.Println("URL: ", r.URL)
		fmt.Println("Method: ", r.Method)
		fmt.Println("Headers: ", r.Header)
		fmt.Println("RemoteAddr: ", r.RemoteAddr)
		fmt.Println("===========================================")
		next.ServeHTTP(w, r)
	})
}
