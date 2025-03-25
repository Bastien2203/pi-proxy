package middlewares

import "net/http"

type Middleware struct {
	Name    string                 `json:"name"`
	Options map[string]interface{} `json:"options"`
}

func ApplyMiddlewares(handler http.Handler, middlewares []Middleware) http.Handler {
	for _, middleware := range middlewares {
		switch middleware.Name {
		case "LogRequest":
			handler = LogRequestMiddleware(handler)
		case "RateLimiter":
			handler = RateLimiterMiddleware(handler, middleware.Options)
		}
	}
	return handler
}
