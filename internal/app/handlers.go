package app

import (
	"encoding/json"
	"net/http"

	"github.com/johnmackenzie91/azlyrics-fetcher"
)

// PostFetch is the /fetch endpoint
func (a App) PostFetch(w http.ResponseWriter, r *http.Request) {
	var err error
	input := azlyrics.FetchRequest{}
	out := azlyrics.FetchResponse{}

	if input, err = azlyrics.NewFetchRequestFromJSON(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, out)
		return
	}

	defer r.Body.Close()

	if err = input.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, out)
		return
	}

	out.Lyrics, err = a.client.GetLyrics(r.Context(), input.Artist, input.Title)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		// TODO add err here
		writeResponse(w, out)
		return
	}

	writeResponse(w, out)
}

func writeResponse(w http.ResponseWriter, output interface{}) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(output)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
