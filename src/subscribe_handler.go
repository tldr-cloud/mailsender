package p

import (
	"fmt"
	"github.com/keighl/postmark"
	"html"
	"log"
	"net/http"
)

func SendWelcomeMail(mail string) error {
	serverApiSecretName := "projects/346496881273/secrets/postmarkapp-server-API-token/versions/latest"
	accountApiSecretName := "projects/346496881273/secrets/postmarkapp-account-API-token/versions/latest"
	serverApiSecret, err := accessSecretVersion( serverApiSecretName)
	if err != nil {
		return err
	}
	accountApiSecret, err := accessSecretVersion( accountApiSecretName)
	if err != nil {
		return err
	}
	client := postmark.NewClient(serverApiSecret, accountApiSecret)

	email := postmark.TemplatedEmail {
		TemplateId: 20032898,
		InlineCss: true,
		TemplateModel: map[string]interface{}{
			"userName": "bobby joe",
		},
		From:      "do-not-reply@over.news",
		To:        mail,
	}

	_, err = client.SendTemplatedEmail(email)
	if err != nil {
		return err
	}
	return nil
}

func ProcessNewSubscribeMsg(w http.ResponseWriter, r *http.Request) {
	request, err := UnpackSubscribeRequest(r)

	if err != nil {
		log.Println("error during the request unpack: ", err.Error())
		return
	}

	if request.Mail == "" {
		log.Println("message is empty")
		return
	}
	fmt.Println("request received for the user: ", request.Mail)

	subscribed, err := AlreadySubscribed(request.Mail)
	if err != nil {
		log.Println("error during the check is the user: ",
			request.Mail, "is subscribed or not: ", err.Error())
		return
	}

	if subscribed {
		fmt.Println("user: ", request.Mail, "is subscribed already")
		return
	}

	err = AddMailToDB(request.Mail)
	if err != nil {
		log.Println("error adding mail: ", request.Mail, " to the DB, error: ", err.Error())
		return
	}
	err = SendWelcomeMail(request.Mail)
	if err != nil {
		log.Println("error adding mail: ", request.Mail, " to the DB, error: ", err.Error())
		return
	}
	fmt.Println("request for the user: ", request.Mail, " been processed correctly")
	fmt.Fprint(w, html.EscapeString("done"))
}
