package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestQueryInt(t *testing.T) {
	tests := []struct {
		name       string
		queryKey   string
		queryValue string
		defaultVal int
		want       int
	}{
		{
			name:       "valid integer",
			queryKey:   "page",
			queryValue: "10",
			defaultVal: 1,
			want:       10,
		},
		{
			name:       "negative integer returns default",
			queryKey:   "page",
			queryValue: "-1",
			defaultVal: 1,
			want:       1,
		},
		{
			name:       "zero returns default",
			queryKey:   "page",
			queryValue: "0",
			defaultVal: 1,
			want:       1,
		},
		{
			name:       "invalid integer returns default",
			queryKey:   "page",
			queryValue: "abc",
			defaultVal: 1,
			want:       1,
		},
		{
			name:       "empty value returns default",
			queryKey:   "page",
			queryValue: "",
			defaultVal: 1,
			want:       1,
		},
		{
			name:       "missing parameter returns default",
			queryKey:   "page",
			queryValue: "", // not set in URL
			defaultVal: 1,
			want:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request with query parameter
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.queryValue != "" {
				q := req.URL.Query()
				q.Add(tt.queryKey, tt.queryValue)
				req.URL.RawQuery = q.Encode()
			}
			c.Request = req

			got := QueryInt(c, tt.queryKey, tt.defaultVal)
			require.Equal(t, tt.want, got)
		})
	}
}
