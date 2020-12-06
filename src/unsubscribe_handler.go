package p

import (
	"fmt"
	"html"
	"log"
	"net/http"
)


func ProcessUnSubscribeMsg(w http.ResponseWriter, r *http.Request) {
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

	err = RemoveMailFromDB(request.Mail)
	if err != nil {
		log.Println("error adding mail: ", request.Mail, " to the DB, error: ", err.Error())
		return
	}
        // Set CORS headers for the preflight request
        if r.Method == http.MethodOptions {
                w.Header().Set("Access-Control-Allow-Origin", "http://over.news/")
                w.Header().Set("Access-Control-Allow-Methods", "POST")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
                w.Header().Set("Access-Control-Max-Age", "3600")
                w.WriteHeader(http.StatusNoContent)
                return
        }
        // Set CORS headers for the main request.
        w.Header().Set("Access-Control-Allow-Origin", "http://over.news/")

	fmt.Fprint(w, html.EscapeString("done"))
}
