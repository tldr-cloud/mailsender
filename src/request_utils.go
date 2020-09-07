package p

import (
	"encoding/json"
	"log"
	"net/http"
)

type request struct{
	Mail string `json:"mail"`
}

func UnpackRequest(r *http.Request) (request, error){
	var request request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println( "error: ", err.Error())
		return request, err
	}
	return request, nil
}
