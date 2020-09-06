package main

import (
	"github.com/keighl/postmark"
)

func main() {
	serverApiSecretName := "projects/346496881273/secrets/postmarkapp-server-API-token/versions/latest"
	accountApiSecretName := "projects/346496881273/secrets/postmarkapp-account-API-token/versions/latest"
	serverApiSecret, err := accessSecretVersion( serverApiSecretName)
	if err != nil {
		panic(err)
	}
	accountApiSecret, err := accessSecretVersion( accountApiSecretName)
	if err != nil {
		panic(err)
	}
	client := postmark.NewClient(serverApiSecret, accountApiSecret)

	email := postmark.Email{
		From:       "do-not-reply@tldr.cloud",
		To:         "viacheslav@kovalevskyi.com",
		Subject:    "test",
		HtmlBody:   "...",
		TextBody:   "...",
		Tag:        "pw-reset",
		TrackOpens: true,
	}

	_, err = client.SendEmail(email)
	if err != nil {
		panic(err)
	}
}
