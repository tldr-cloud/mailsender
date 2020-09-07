package p

import (
	"fmt"
	"html"
	"log"
	"net/http"
)


func ProcessUnSubscribeMsg(w http.ResponseWriter, r *http.Request) {
	request, err := UnpackRequest(r)

	if err != nil {
		log.Println("error during the request unpack: ", err.Error())
		return
	}

	if request.Mail == "" {
		log.Println("message is empty")
		return
	}
	fmt.Println("request received for the user: ", request.Mail)

	err = RemoveMailFromDB(request.Mail)
	if err != nil {
		log.Println("error adding mail: ", request.Mail, " to the DB, error: ", err.Error())
		return
	}
	fmt.Fprint(w, html.EscapeString("done"))
}
