package feedmgmt_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/feedmgmt"
	"github.com/stretchr/testify/require"
)

func TestGetFeeds(t *testing.T) {
	tests := map[string]struct {
		provider           string
		category           string
		enabled            bool
		returnedStatusCode int
		expectedError      bool
	}{
		"test 1": {
			returnedStatusCode: 500,
			expectedError:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/v1/feeds", r.URL.Path)
				require.Equal(t, "GET", r.Method)
				w.WriteHeader(test.returnedStatusCode)
			}))
			defer ts.Close()

			client := feedmgmt.NewClient(ts.URL, ts.Client())

			feeds, err := client.GetFeeds("", "", false)

			_ = feeds

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
