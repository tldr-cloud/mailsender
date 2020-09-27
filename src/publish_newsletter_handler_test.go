package p

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConvertTldrToHtml(t *testing.T) {
	expectedHtml := `<div>
    <a href="http://techcrunch.com/2017/02/23/website-builder-wix-acquires-art-community-deviantart-for-36m/">Website builder Wix acquires art community DeviantArt for $36M – TechCrunch</a>
    <img src="https://techcrunch.com/wp-content/uploads/2017/02/screen-shot-2017-02-23-at-12-59-34.png?w=730">
    <div>
        Wix .com has made another acquisition to build out the tools that it provides to users to build and administer websites: it has acquired DeviantArt, an online community for artists, designers and art/design enthusiasts with some 325 million individual pieces of original art and more than 40 million registered members, for $36 million in cash, including $3 million of assumed liabilities. Updated detail related to DeviantArt’s valuation prior to its sale.
    </div>
</div>
`

	html, err := ConvertTldrToHtml(
		"http___techcrunch.com_2017_02_23_website-" +
			"builder-wix-acquires-art-community-deviantart-for-36m_")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, expectedHtml, string(html))
}

func TestConvertNewsletterToHtml(t *testing.T) {
	expectedHtml := ``

	html, err := ConvertNewsletterToHtml("20200927161650")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, expectedHtml, string(html))
}
// 20200927161650
