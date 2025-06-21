package main

import (
	"net/http"
	"fmt"
)
	
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	c := fmt.Sprintf(`
		<html>
  	<body>
    	<h1>Welcome, Chirpy Admin</h1>
    	<p>Chirpy has been visited %d times!</p>
  	</body>
		</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(c))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
