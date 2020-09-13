package p

import (
	"context"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func PublishNewsletter(ctx context.Context, m PubSubMessage) error {
	newsletterId := string(m.Data)
	newletter, err := GetNewsletterById(newsletterId)
	if err != nil {
		return err
	}
	for _, tldrId := range newletter.NewsIds {

	}
	return nil
}
