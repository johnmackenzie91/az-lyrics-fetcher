package app_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"net/http/httptest"

	"context"

	"github.com/johnmackenzie91/azlyrics-fetcher"
	"github.com/johnmackenzie91/azlyrics-fetcher/internal/app"
	"github.com/johnmackenzie91/azlyrics-fetcher/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

var stubCtx = context.Background()

//go:generate mockery -name=LyricClient
func TestApp_PostFetch(t *testing.T) {
	requestBody, err := json.Marshal(azlyrics.FetchRequest{
		Artist: "some artist",
		Title:  "some song title",
	})
	assert.Nil(t, err)

	mockClient := mocks.LyricClient{}
	mockClient.On("GetLyrics", stubCtx, "some artist", "some song title").
		Return("some lyrics ...", nil)

	r, err := http.NewRequestWithContext(stubCtx, "GET", "http://example.com", bytes.NewBuffer(requestBody))
	assert.Nil(t, err)

	sut, _ := app.New(&mockClient)

	w := httptest.NewRecorder()
	sut.PostFetch(w, r)

	out := azlyrics.FetchResponse{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&out))

	expected := azlyrics.FetchResponse{
		Lyrics: "some lyrics ...",
	}
	assert.Equal(t, expected, out)
	mockClient.AssertExpectations(t)
}
