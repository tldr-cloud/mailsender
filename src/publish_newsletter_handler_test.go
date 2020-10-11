package p

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"io/ioutil"
)

func TestConvertNewsletterToHtml(t *testing.T) {
	expectedHtml, err := ioutil.ReadFile("templates/index-test.html")
	if err != nil {
		t.Error(err.Error())
		return
	}

	html, err := ConvertNewsletterToHtml("20200927161650")
	fmt.Printf("generated html: \n%s\n", html)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, expectedHtml, html)
}
