package client

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"errors"

	"github.com/PuerkitoBio/goquery"
)

type doer interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client will be used to make requests to azlyrics.
type Client struct {
	client   doer
	endpoint *url.URL
}

var defaultEndpoint = mustParse("https://search.azlyrics.com/")

// New constructs a new client
func New(ops ...Option) (Client, error) {
	// sensible defaults
	c := Client{
		// this is the correct domain to hit
		endpoint: defaultEndpoint,
		// basic default go client.
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, op := range ops {
		if err := op(&c); err != nil {
			return Client{}, err
		}
	}

	return c, nil
}

// GetLyrics attempts to find the song lyrics on azlyrics
func (c Client) GetLyrics(ctx context.Context, artistName, songTitle string) (string, error) {
	u, err := c.searchPage(ctx, artistName, songTitle)
	if err != nil {
		return "", err
	}

	return c.lyricsPage(ctx, u)
}

// search requests and parses azlyrics search page
func (c Client) searchPage(ctx context.Context, artistName, songTitle string) (string, error) {
	res, err := c.getSearchPage(ctx, artistName, songTitle)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return "", err
	}

	u, err := c.parseSearchPage(doc, artistName, songTitle)
	return u, nil
}

// getSearchPage fetches the search page by searching for artist name and songTitle
func (c Client) getSearchPage(ctx context.Context, artistName, songTitle string) (*http.Response, error) {
	query := "/search.php?q=" + url.QueryEscape(artistName+" "+songTitle)
	domain := c.endpoint.Scheme + "://" + c.endpoint.Host
	r, err := http.NewRequestWithContext(ctx, "GET", domain+query, nil)

	if err != nil {
		return nil, err
	}

	return c.client.Do(r)
}

// parseSearchPage parses the search page for the song specified
func (c Client) parseSearchPage(doc *goquery.Document, artistName, songTitle string) (string, error) {
	var outURL string
	doc.Find("table tr").Each(func(i int, tr *goquery.Selection) {
		if outURL != "" {
			return
		}
		var title string
		var artist string
		var pageURL string
		if anchor := tr.Find("a"); anchor != nil {
			v, ok := anchor.Attr("href")

			if !ok {
				return
			}
			pageURL = v
			title = strings.Replace(anchor.Text(), "\"", "", 2)
			title = strings.ToLower(title)
			title = strings.TrimSpace(title)
		}
		if bold := tr.Find("b:nth-child(2)"); bold != nil {
			artist = bold.Text()
			artist = strings.ToLower(artist)
			artist = strings.TrimSpace(artist)
		}

		if artistName == artist && songTitle == title {
			outURL = pageURL
		}
		return
	})
	return outURL, nil
}

func (c Client) lyricsPage(ctx context.Context, u string) (string, error) {
	r, err := http.NewRequestWithContext(ctx, "GET", u, nil)

	if err != nil {
		return "", err
	}

	res, err := c.client.Do(r)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return "", err
	}

	return c.parseLyricPage(doc)
}

var missingBrRegex = regexp.MustCompile(`(([a-zA-Z0-9]{1})([A-Z]))`)

// parseSearchPage parses the search page for the song specified
func (c Client) parseLyricPage(doc *goquery.Document) (string, error) {
	lyricDiv := doc.Find("div:nth-of-type(5)").First()
	if lyricDiv == nil {
		return "", errors.New("unable to find lyrics")
	}

	out := strings.Replace(lyricDiv.Text(), "\n", "", -1)
	out = strings.Replace(out, "?", " ", -1)
	out = missingBrRegex.ReplaceAllString(out, "$2 $3")

	return out, nil
}

func mustParse(u string) *url.URL {
	out, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return out
}
