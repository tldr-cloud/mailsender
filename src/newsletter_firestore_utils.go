package p

import (
	"cloud.google.com/go/firestore"
	"context"
)

var newsletters *firestore.CollectionRef

type news struct {
	Title string  `json:"title"`
	PicUrl string  `json:"pic_url"`
	Text string  `json:"text"`
	Url string  `json:"url"`
}

type newsletter struct {
	News []news `json:"news"`
}

func MaybeInitNewslettersCollection() error {
	ctx := context.Background()
	if newsletters == nil {
		firestoreClient, err := firestore.NewClient(ctx, projectId)
		if err != nil {
			return err
		}

		newsletters = firestoreClient.Collection("newsletters")
	}
	return nil
}

func GetListOfNewsForNewsletter(newsletterId string) (newsletter, error) {
	if err := MaybeInitNewslettersCollection(); err != nil {
		return newsletter{}, err
	}
	// TODO
	return newsletter{}, nil
}
