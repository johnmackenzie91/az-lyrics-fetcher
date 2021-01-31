package azlyrics

import (
"encoding/json"
"io"
"strings"
)

func NewFetchRequestFromJSON(r io.ReadCloser) (FetchRequest, error) {
	out := FetchRequest{}
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		return out, Error{
			Code:    EINVALID,
			Message: "failed to json decode request body",
		}
	}
	return out, nil
}

// UnmarshalJSON unmarshals input and converts all values to lowercase
func (f *FetchRequest) UnmarshalJSON(b []byte) error {
	aux := struct {
		Artist string `json:"artist"`
		Title  string `json:"title"`
	}{}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	f.Artist = strings.ToLower(aux.Artist)
	f.Title = strings.ToLower(aux.Title)

	return nil
}

var errMissingRequestField = Error{
	Code:    EINVALID,
	Message: "field missing from request body. Fields needed; title & artist",
}

// Validate valides the user input by checking that all fields on the struct are set
func (f FetchRequest) Validate() error {
	if f.Title == "" || f.Artist == "" {
		return errMissingRequestField
	}
	return nil
}
