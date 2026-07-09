package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		header string
		want   string
	}{
		{"Bearer lnk_abc", "lnk_abc"},
		{"bearer lnk_abc", "lnk_abc"},
		{"BEARER lnk_abc", "lnk_abc"},
		{"Bearer   lnk_abc  ", "lnk_abc"},
		{"", ""},
		{"lnk_abc", ""},
		{"Basic lnk_abc", ""},
		{"Bearer", ""},
		{"Token lnk_abc", ""},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		if tc.header != "" {
			c.Request.Header.Set("Authorization", tc.header)
		}
		if got := bearerToken(c); got != tc.want {
			t.Fatalf("bearerToken(%q) = %q, want %q", tc.header, got, tc.want)
		}
	}
}
