package p

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}


func SendNewsletter(newsletterHtml string) error {
	toMails := make([]*mail.Email, 0)
	toMailsStrings, err := GetMailLists()
	if err != nil {
		return err
	}
	for _, mailString := range toMailsStrings {
		toMails = append(toMails, mail.NewEmail("", mailString))
	}
	sendGridApiKey, err := accessSecretVersion("sendgrid-com-newsletter-apikey")
	if err != nil {
		return err
	}
	now := time.Now()
	from := mail.NewEmail("TLDR Newsletter", "newsletter@tldr.cloud")
	subject := fmt.Sprintf("newsletter for: %s", now.String())

	mailsToSend := make([]*mail.SGMailV3, len(toMails))
	for _, toMail := range toMails {
		mailsToSend = append(mailsToSend,
			mail.NewSingleEmail(from, subject, toMail, "", newsletterHtml))
	}

	client := sendgrid.NewSendClient(sendGridApiKey)
	for _, mailToSend := range mailsToSend {
		client.Send(mailToSend)
	}
 	return nil
}

func PubSubMessageHandler(ctx context.Context, m PubSubMessage) error {
	newsletterId := string(m.Data)
	fmt.Printf("news letter with id: %s received\n", newsletterId)
	if err := PublishNewsletter(newsletterId); err != nil {
		log.Printf("main error for the id: %s is %s\n", newsletterId, err.Error())
		return err
	}
	return nil
}


func ConvertTldrToHtml(tldrId string) (template.HTML, error) {
	newsTemplate, err := template.ParseFiles("templates/news.gohtml")
	if err != nil {
		log.Printf("new.gohtml can't be parsed due to the error: %s\n",
			err.Error())
		return "", err
	}

	var buf bytes.Buffer
	tldr, err := GetTldrById(tldrId)
	if err != nil {
		log.Printf("tldr with id: %s can't be converted to HTML due to error: %s",
			tldrId, err.Error())
		return "", err
	}
	if err = newsTemplate.Execute(&buf, tldr); err != nil {
		log.Printf("hatpm teplate can' be applied to tldr with id: %s due to error: %s\n",
			tldrId, err.Error())
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func ConvertNewsletterToHtml(newsletterId string) (string, error) {
	fmt.Printf("starting convertions for the newsletter with id: %s", newsletterId)
	newsletter, err := GetNewsletterById(newsletterId)
	if err != nil {
		log.Printf("error getting newsletter (id: %s) from the db. Error: %s\n",
			newsletterId, err.Error())
		return "", err
	}
	fmt.Printf("for news with id: %s amount of tldrs is: %d", newsletterId, len(newsletter.NewsIds))
	records := make([]template.HTML, len(newsletter.NewsIds))
	for index, tldrId := range newsletter.NewsIds {
		if records[index], err = ConvertTldrToHtml(tldrId); err != nil {
			return "", err
		}
		fmt.Printf("news with index %d got converted to: %s\n", index, records[index])
	}

	newsletterTemplate, err := template.ParseFiles("templates/newsletter.gohtml")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	newsletterTemplate.Execute(&buf, records)
	return buf.String(), nil
}

func PublishNewsletter(newsletterId string) error {
	fmt.Printf("news letter with id: %s received", newsletterId)
	html, err := ConvertNewsletterToHtml(newsletterId)
	if err != nil {
		return err
	}
	SendNewsletter(html)
	return nil
}
