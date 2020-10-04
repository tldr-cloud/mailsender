package p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertNewsletterToHtml(t *testing.T) {
	expectedHtml := `<html>
    <body>
    
        <a href="http://techcrunch.com/2017/02/23/website-builder-wix-acquires-art-community-deviantart-for-36m/">Website builder Wix acquires art community DeviantArt for $36M – TechCrunch</a>
        <img src="https://techcrunch.com/wp-content/uploads/2017/02/screen-shot-2017-02-23-at-12-59-34.png?w=730">
        <div>
            Wix .com has made another acquisition to build out the tools that it provides to users to build and administer websites: it has acquired DeviantArt, an online community for artists, designers and art/design enthusiasts with some 325 million individual pieces of original art and more than 40 million registered members, for $36 million in cash, including $3 million of assumed liabilities. Updated detail related to DeviantArt’s valuation prior to its sale.
        </div>
    
        <a href="http://techcrunch.com/2020/07/02/reliance-jio-platforms-launches-jiomeet-video-conference-service/">India’s richest man takes on Zoom – TechCrunch</a>
        <img src="https://techcrunch.com/wp-content/uploads/2020/07/GettyImages-1203053791.jpg?w=586">
        <div>
            India’s Reliance Jio Platforms, which recently concluded a $15.2 billion fundraise run, is ready to enter a new business: Video conferencing. Like Zoom and Google Meet, JioMeet offers unlimited number of free calls in high definition (720p) to users and supports as many as 100 participants on a call. JioMeet is available for use through Chrome and Firefox browsers on desktop, as well as has standalone apps for macOS, Windows, iOS, and Android.
        </div>
    
    </body>
</html>`

	html, err := ConvertNewsletterToHtml("20200927161650")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, expectedHtml, string(html))
}
