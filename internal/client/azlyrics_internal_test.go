package client

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/johnmackenzie91/httptestfixtures"
	"github.com/stretchr/testify/assert"
)

func TestClient_parseSearchPage(t *testing.T) {
	// arrange mock input
	page := httptestfixtures.MustLoadRequest(t, "./testdata/search_arctic_monkeys_505")
	doc, err := goquery.NewDocumentFromReader(page.Body)
	assert.Nil(t, err)
	assert.Nil(t, page.Body.Close())

	sut := Client{}
	out, err := sut.parseSearchPage(doc, "arctic monkeys", "505")
	assert.Nil(t, err)

	expected := "https://www.azlyrics.com/lyrics/arcticmonkeys/505.html"
	assert.Equal(t, expected, out)
}

//go:generate curl -i -L -A "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0" https://www.azlyrics.com/lyrics/arcticmonkeys/505.html -o ./testdata/arctic_monkeys_505
//go:generate mockery -name=doer -outpkg=client_test -testonly -output=.
func TestClient_parseLyricPage(t *testing.T) {
	// arrange mock input
	page := httptestfixtures.MustLoadRequest(t, "./testdata/arctic_monkeys_505")
	doc, err := goquery.NewDocumentFromReader(page.Body)
	assert.Nil(t, err)
	assert.Nil(t, page.Body.Close())

	sut := Client{}
	out, err := sut.parseLyricPage(doc)
	assert.Nil(t, err)

	expected := "I'm going back to 505 If it's a 7 hour flight or a 45 minute drive In my imagination you're waiting lying on your side With your hands between your thighs Stop and wait a sec Oh when you look at me like that my darling What did you expect I probably still adore you with your hands around my neck Or I did last time I checked Not shy of a spark A knife twists at the thought that I should fall short of the mark Frightened by the bite though its no harsher than the bark Middle of adventure, such a perfect place to start I'm going back to 505 If it's a 7 hour flight or a 45 minute drive In my imagination you're waiting lying on your side With your hands between your thighs But I crumble completely when you cry It seems like once again you've had to greet me with goodbye I'm always just about to go and spoil a surprise Take my hands off of your eyes too soon I'm going back to 505 If it's a 7 hour flight or a 45 minute drive In my imagination you're waiting lying on your side With your hands between your thighs...and a smile"
	assert.Equal(t, expected, out)
}

//go:generate curl -i -L -A "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0" https://www.azlyrics.com/lyrics/killers/whenyouwereyoung.html -o ./testdata/the_killers_when_you_were_young
func TestClient_parseLyricPage_When_You_were_Young(t *testing.T) {
	// arrange mock input
	page := httptestfixtures.MustLoadRequest(t, "./testdata/the_killers_when_you_were_young")
	doc, err := goquery.NewDocumentFromReader(page.Body)
	assert.Nil(t, err)
	assert.Nil(t, page.Body.Close())

	sut := Client{}
	out, err := sut.parseLyricPage(doc)
	assert.Nil(t, err)

	expected := "You sit there in your heartache Waiting on some beautiful boy to To save you from your old ways You play forgiveness Watch it now Here he comes He doesn't look a thing like Jesus But he talks like a gentleman Like you imagined When you were young Can we climb this mountain I don't know Higher now than ever before I know we can make it if we take it slow Let's take it easy Easy now Watch it go We're burning down the highway skyline On the back of a hurricane That started turning When you were young When you were young And sometimes you close your eyes And see the place where you used to live When you were young They say the devil's water It ain't so sweet You don't have to drink right now But you can dip your feet Every once in a little while You sit there in your heartache Waiting on some beautiful boy to To save you from your old ways You play forgiveness Watch it now Here he comes He doesn't look a thing like Jesus But he talks like a gentleman Like you imagined When you were young(Talks like a gentleman)(Like you imagined)When you were young I said he doesn't look a thing like Jesus He doesn't look a thing like Jesus But more than you'll ever know"
	assert.Equal(t, expected, out)
}
