package p

import (
	"log"
	"net/http"
)

func ProcessNewSubscribeConfirmationMsg(w http.ResponseWriter, r *http.Request) {
	codeFromRequest := GetSubscriptionConfirmationCodeFromRequest(r)
	mail := GetSubscriptionMailFromQueryFromRequest(r)
	codeFromDb, err := GetMailVerificationCodeFromDb(mail)
	if err != nil {
		log.Println("error during getting the code for mail: ", mail)
		log.Println("error: ", err.Error())
		return
	}

	if codeFromDb == codeFromRequest {
		err = MarkMailAddressAsVerified(mail)
		if err != nil {
			log.Printf("error marking address %s as verified: %s\n", mail, err.Error())
		}
		return
	}

	log.Printf("for mail: %s, code from DB (%s) does not match code from the request (%s)\n",
		mail, codeFromDb, codeFromRequest)
	return


}
