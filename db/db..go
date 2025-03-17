package db

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
)

const projectID = "enhanced-rarity-437111-d9"

type Store struct {
	db *firestore.Client
}

type User struct {
	ChatId           int64
	UserName         string
	SentenceLanguage string //Language in which sentence should be generated
	Level            string //e.g. A1
	PremiumUntil     int64  //unix time
	State            int
}

func New(ctx context.Context) (*Store, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Store{client}, nil
}

func (store *Store) SaveUser(ctx context.Context, user *User) error {
	_, err := store.db.Collection("users").Doc(strconv.Itoa(int(user.ChatId))).Set(ctx, user)
	return err
}

func (store *Store) GetUser(ctx context.Context, chatId string) (*User, error) {
	res, err := store.db.Collection("users").Doc(chatId).Get(ctx)
	if err != nil {
		return nil, err
	}
	data := res.Data()
	return &User{
		ChatId:           data["ChatId"].(int64),
		UserName:         data["UserName"].(string),
		SentenceLanguage: data["SentenceLanguage"].(string),
		Level:            data["Level"].(string),
		PremiumUntil:     data["PremiumUntil"].(int64),
		State:            int(data["State"].(int64)),
	}, nil
}

func (store *Store) UpdateUserState(ctx context.Context, chatId string, state int) error {
	_, err := store.db.Collection("users").Doc(chatId).Update(ctx, []firestore.Update{
		{
			Path:  "State",
			Value: state,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) UpdateUserPremium(ctx context.Context, chatId string, premiumUntil int64) error {
	_, err := store.db.Collection("users").Doc(chatId).Update(ctx, []firestore.Update{
		{
			Path:  "PremiumUntil",
			Value: premiumUntil,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) UpdateUserLevel(ctx context.Context, chatId string, level string) error {
	_, err := store.db.Collection("users").Doc(chatId).Update(ctx, []firestore.Update{
		{
			Path:  "Level",
			Value: level,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) DeleteUser(ctx context.Context, chatId string) error {
	_, err := store.db.Collection("users").Doc(chatId).Delete(ctx)
	return err
}
