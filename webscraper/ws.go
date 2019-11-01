package webscraper

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// DownloadLink is the type to aggregate the download data
type DownloadLink struct {
	filename string
	href     string
	sha256   string
}

// FileName returns the filename value of the download link
func (dl DownloadLink) FileName() string {
	return dl.filename
}

// Href returns the href value of the download link
func (dl DownloadLink) Href() string {
	return dl.href
}

// Sha256 returns the Sha256 value of the download link
func (dl DownloadLink) Sha256() string {
	return dl.sha256
}

// GetHighlightClassTokensN will return the up to the first N of all the highligth class tokens in the
// http.Response along with the sha value that is on the table
func GetHighlightClassTokensN(resp *http.Response, n int) ([]DownloadLink, error) {

	links := make([]DownloadLink, 0, n)

	z := html.NewTokenizer(resp.Body)

	for {
		tokenType := z.Next()

		switch {
		case tokenType == html.ErrorToken:
			// End of the document, we're done
			return links, nil
		case tokenType == html.StartTagToken:
			token := z.Token()

			if isTr(&token) && isClassHighlight(&token) {
				anchorToken := getNextAnchor(z)
				if anchorToken == nil {
					return links, nil
				}
				href, ok := getHTTPHref(anchorToken)
				if !ok {
					continue
				}
				z.Next()
				fileName := z.Token().Data
				getNextTT(z)
				z.Next()
				sha256 := z.Token().Data

				links = append(links, DownloadLink{fileName, href, sha256})
				if len(links) == n {
					return links, nil
				}
			}
		}
	}
}

// getHTTPHref will return the token's href value if starts with "http" otherwise ok will be false
func getHTTPHref(token *html.Token) (string, bool) {
	href, ok := getHref(token)
	if !ok {
		return "", false
	}
	// Make sure the url begins in http**
	hasProto := strings.Index(href, "http") == 0
	if hasProto {
		return href, true
	}
	return "", false
}

// getNextTT will iterate on the tokenizer until it finds a token type <tt> or the end of the document is reached
func getNextTT(z *html.Tokenizer) *html.Token {
	for {
		tokenType := z.Next()
		switch {
		case tokenType == html.ErrorToken:
			return nil
		case tokenType == html.StartTagToken:
			token := z.Token()
			if isTT(&token) {
				return &token
			}
		}
	}
}

// getNextAnchor will iterate on the tokenizer until it finds a token type <a> or the end of the document is reached
func getNextAnchor(z *html.Tokenizer) *html.Token {
	for {
		tokenType := z.Next()
		switch {
		case tokenType == html.ErrorToken:
			return nil
		case tokenType == html.StartTagToken:
			token := z.Token()
			if isAnchor(&token) {
				return &token
			}
		}
	}
}

// isTr returns true if the specific Token is <tr>
func isTr(token *html.Token) bool {
	return token.Data == "tr"
}

// isAnchor returns true if the specific Token is an anchor <a>
func isAnchor(token *html.Token) bool {
	return token.Data == "a"
}

// isTT returns true if the specific Token is <tt>
func isTT(token *html.Token) bool {
	return token.Data == "tt"
}

// isClassHighlight returns true if the specific Token has a class attribute
// with the value "highlight"
func isClassHighlight(token *html.Token) bool {
	class, _ := getClass(token)
	if class == "highlight" {
		return true
	}
	return false
}

// Helper function to pull the class attribute from a Token
func getClass(t *html.Token) (string, bool) {
	// Iterate over all of the Token's attributes until we find a "class"
	for _, attr := range t.Attr {
		if attr.Key == "class" {
			return attr.Val, true
		}
	}
	return "", false
}

// Helper function to pull the href attribute from a Token
func getHref(t *html.Token) (string, bool) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			return attr.Val, true
		}
	}
	return "", false
}
