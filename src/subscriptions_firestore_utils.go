package p

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"math/rand"
	"strconv"
	"time"
)

var subscribers *firestore.CollectionRef

func MaybeInit() error {
	ctx := context.Background()
	if subscribers == nil {
		firestoreClient, err := firestore.NewClient(ctx, mailSenderProjectId)
		if err != nil {
			return err
		}

		subscribers = firestoreClient.Collection("subscribers")
	}
	return nil
}

func GetMailLists() ([]string, error) {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return nil, err
	}
	iter := subscribers.Documents(ctx)

	mails := make([]string, 0)

	for {
		subscriber, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		mail, err := subscriber.DataAt("email")
		if err != nil {
			return nil, err
		}
		mails = append(mails, mail.(string))
	}
	return mails, nil
}

func AlreadySubscribed(mail string) (bool, error) {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return false, err
	}
	iter := subscribers.Where("email", "==", mail).Documents(ctx)
	_, err = iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

func GenerateCustomerConfirmationId() string {
	const customerIdSize = 10
	id := ""
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < customerIdSize; i++ {
		id = id + strconv.Itoa(seededRand.Intn(10))
	}
	return id
}

func AddMailToDB(mail string) error {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return err
	}
	_, _, err = subscribers.Add(ctx, map[string]interface{}{
		"email":    mail,
		"subscribedDate": time.Now(),
		"verified": false,
		"verificationCode": GenerateCustomerConfirmationId(),
	})
	if err != nil {
		return err
	}
	return nil
}

func RemoveMailFromDB(mail string) error {
	ctx := context.Background()
	MaybeInit()
	subscribed, err := AlreadySubscribed(mail)
	if err != nil {
		return err
	}
	if !subscribed {
		return nil
	}
	iter := subscribers.Where("email", "==", mail).Documents(ctx)
	record, err := iter.Next()
	if err == iterator.Done {
		return nil
	}
	_, err = record.Ref.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetMailVerificationCodeFromDb(mail string) (string, error) {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return "", err
	}
	iter := subscribers.Where("email", "==", mail).Documents(ctx)
	subscriber, err := iter.Next()
	if err == iterator.Done {
		return "", errors.New("no mail found")
	}
	if err != nil {
		return "", err
	}
	code, err := subscriber.DataAt("verificationCode")
	if err != nil {
		return "", err
	}
	return code.(string), nil
}

func MarkMailAddressAsVerified(mail string) error {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return err
	}
	iter := subscribers.Where("email", "==", mail).Limit(1).Documents(ctx)
	subscriber, err := iter.Next()
	if err == iterator.Done {
		return errors.New("no mail found")
	}
	_, err = subscribers.Doc(subscriber.Ref.ID).Set(ctx,
		map[string]interface{}{
		"verified": true,
	})
	if err != nil {
		return nil
	}
	return nil
}
