package p

import (
	"bytes"
	"context"
	"html/template"
	"github.com/keighl/postmark"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}


func SendNewsletter(newsletterHtml string) error {
	mails := GetMailLists()

	// TODO
	return nil
}


func PublishNewsletter(ctx context.Context, m PubSubMessage) error {
	newsletterId := string(m.Data)
	newletter, err := GetNewsletterById(newsletterId)
	if err != nil {
		return err
	}
	newsTemplate, err := template.ParseFiles("templates/news.gohtml")
	if err != nil {
		return err
	}
	resultRecords := make([]string, len(newsletterId))
	for index, tldrId := range newletter.NewsIds {
		var buf bytes.Buffer
		tldr, err := GetNewsById(tldrId)
		if err != nil {
			return err
		}
		newsTemplate.Execute(&buf, tldr)
		resultRecords[index] = buf.String()
	}

	newsletterTemplate, err := template.ParseFiles("templates/newsletter.gohtml")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	newsletterTemplate.Execute(&buf, resultRecords)
	newsletterHtml := buf.String()
	SendNewsletter(newsletterHtml)
	return nil
}
