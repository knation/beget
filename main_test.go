package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Makes an HTTP test to the given `url` with the given `body`
func httpTest(r *gin.Engine, method string, url string, body io.Reader) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	r.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

// Makes an HTTP tes request using the given `req`
// func httpTestRequest(r *gin.Engine, req *http.Request) (int, string) {
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)
// 	return w.Code, w.Body.String()
// }

func TestHealthCheck(t *testing.T) {
	r := initRouter()

	code, body := httpTest(r, http.MethodGet, "/healthz", nil)

	assert.Equal(t, 200, code)
	assert.Equal(t, "OK", body)
}
