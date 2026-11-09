package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newRouter(apiKey string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", APIKeyAuth(apiKey), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestAPIKeyAuth(t *testing.T) {
	cases := []struct {
		name   string
		key    string
		header string
		want   int
	}{
		{"missing header", "secret", "", http.StatusUnauthorized},
		{"wrong key", "secret", "Bearer nope", http.StatusUnauthorized},
		{"no bearer prefix", "secret", "secret", http.StatusUnauthorized},
		{"empty configured key", "", "Bearer ", http.StatusUnauthorized},
		{"valid key", "secret", "Bearer secret", http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()
			newRouter(tc.key).ServeHTTP(w, req)
			if w.Code != tc.want {
				t.Errorf("got %d, want %d", w.Code, tc.want)
			}
		})
	}
}
