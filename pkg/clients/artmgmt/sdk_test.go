package artmgmt_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/artmgmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddArticles(t *testing.T) {
	tests := map[string]struct {
		articles           artmgmt.Articles
		returnedStatusCode int
		expectedError      bool
	}{
		"test 1": {
			articles: artmgmt.Articles{
				artmgmt.Article{
					GUID: "guid1",
				},
			},
			returnedStatusCode: 500,
			expectedError:      true,
		},
		"test 2": {
			articles: artmgmt.Articles{
				artmgmt.Article{
					GUID: "guid1",
				},
			},
			returnedStatusCode: 204,
			expectedError:      false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/v1/articles:batch", r.URL.Path)
				require.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(test.returnedStatusCode)
			}))
			defer ts.Close()

			client := artmgmt.NewClient(ts.URL, ts.Client())

			err := client.AddArticles(test.articles)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
