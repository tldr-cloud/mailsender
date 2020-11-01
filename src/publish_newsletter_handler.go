package p

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"log"
	"time"
	"os"
)

const newsletterDateLayout = "January 2, 2006"

const gcloudFuncSourceDir = "serverless_function_source_code"

func FixDir() {
	fileInfo, err := os.Stat(gcloudFuncSourceDir)
	if err == nil && fileInfo.IsDir() {
		_ = os.Chdir(gcloudFuncSourceDir)
	}
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func SendNewsletter(newsletterHtml string) error {
	fmt.Println("Sending letters")
	toMails := make([]*mail.Email, 0)
	toMailsStrings, err := GetMailLists()
	if err != nil {
		return err
	}
	for _, mailString := range toMailsStrings {
		fmt.Printf("mail string from the DB: %s\n", mailString)
		toMails = append(toMails, mail.NewEmail("", mailString))
	}
	sendGridApiKey, err := accessSecretVersion(sengridComNewsletterApiKeyName)
	if err != nil {
		return err
	}
	now := time.Now()
	from := mail.NewEmail("TLDR Newsletter", "newsletter@tldr.cloud")
	subject := fmt.Sprintf("newsletter for: %s", now.Format(newsletterDateLayout))

	mailsToSend := make([]*mail.SGMailV3, 0)
	for _, toMail := range toMails {
		fmt.Printf("Mail to send to from the DB: %s\n", toMail)
		mailsToSend = append(mailsToSend,
			mail.NewSingleEmail(from, subject, toMail, "Your daily update is here!", newsletterHtml))
	}

	client := sendgrid.NewSendClient(sendGridApiKey)
	for number, mailToSend := range mailsToSend {
		fmt.Printf("Sending mail number: %d\n", number)
		resp, err := client.Send(mailToSend)
		if err != nil {
			return err
		}
		fmt.Printf("resp body from the sendgrid: %s\n", resp.Body)
		fmt.Printf("resp StatusCode from the sendgrid: %d\n", resp.StatusCode)
	}
	return nil
}

func PubSubMessageHandler(ctx context.Context, m PubSubMessage) error {
	newsletterId := string(m.Data)
	fmt.Printf("newsletter with id: %s received\n", newsletterId)
	if err := PublishNewsletter(newsletterId); err != nil {
		log.Printf("main error for the id: %s is %s\n", newsletterId, err.Error())
		return err
	}
	return nil
}

func ConvertNewsletterToHtml(newsletterId string) (string, error) {
	fmt.Printf("starting convertions for the newsletter with id: %s", newsletterId)
	newsletter, err := GetNewsletterById(newsletterId)
	if err != nil {
		log.Printf("error getting newsletter (id: %s) from the db. Error: %s\n",
			newsletterId, err.Error())
		return "", err
	}
	fmt.Printf("for news with id: %s amount of tldrs is: %d\n", newsletterId, len(newsletter.NewsIds))
	if len(newsletter.NewsIds) == 0 {
		return "", errors.New(fmt.Sprintf("amount of TLDRs in newsleeter (id: %s) is 0", newsletterId))
	}
	records := make([]TLDR, len(newsletter.NewsIds))
	for index, tldrId := range newsletter.NewsIds {
		if records[index], err = GetTldrById(tldrId); err != nil {
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
	FixDir()
	fmt.Printf("newsletter with id: %s received\n", newsletterId)
	html, err := ConvertNewsletterToHtml(newsletterId)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	if err = SendNewsletter(html); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	return nil
}
