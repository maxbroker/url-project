package redirect_test

import (
	"awesomeProject/internal/lib/api"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"awesomeProject/internal/http-server/handlers/redirect"
	"awesomeProject/internal/http-server/handlers/redirect/mocks"
	"awesomeProject/internal/lib/logger/handlers/slogdiscard"
)

func TestRedirectURLHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:  "Success1",
			alias: "test_alias",
			url:   "/fffff",
		},
		{
			name:  "Success2",
			alias: "test_alias",
			url:   "https://www.google1.com/",
		},
		{
			name:  "Success3",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:  "Success4",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:  "Success5",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewUrlGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.RedirectUrlHandler(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// Check the final URL after redirection.
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
