package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"strconv"
)

const projectID = "enhanced-rarity-437111-d9" //Project id on Google cloud

type Store struct {
	db *firestore.Client
}

type User struct {
	ChatId           int64
	UserName         string //Telegram username
	SentenceLanguage string //Language in which sentence should be generated
	Level            string //e.g. A1
	PremiumUntil     int64  //unix time
	PreferencesSet   bool
	LastUsed         int64 //unix time
	FreeSentences    int   //how many more free sentences can user generate
}

// New creates new firestore instance
func New(ctx context.Context) (*Store, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Store{client}, nil
}

// CreateUser Creates user if user does not exist
func (store *Store) CreateUser(ctx context.Context, user *User) error {
	_, err := store.db.Collection("users").Doc(strconv.Itoa(int(user.ChatId))).Create(ctx, user)
	return err
}

// UpdateUser updates user overriding all fields with the provided user struct
func (store *Store) UpdateUser(ctx context.Context, user *User) error {
	_, err := store.db.Collection("users").Doc(strconv.Itoa(int(user.ChatId))).Set(ctx, user)
	return err
}

// GetUser retrieves user from the database using telegram chat id
func (store *Store) GetUser(ctx context.Context, chatId int64) (*User, error) {
	//Make a request to the firestore
	res, err := store.db.Collection("users").Doc(strconv.Itoa(int(chatId))).Get(ctx)
	if err != nil {
		return nil, err
	}

	//Get data from the response
	data := res.Data()

	//Return data in user struct
	return &User{
		ChatId:           data["ChatId"].(int64),
		UserName:         data["UserName"].(string),
		SentenceLanguage: data["SentenceLanguage"].(string),
		Level:            data["Level"].(string),
		PremiumUntil:     data["PremiumUntil"].(int64),
		PreferencesSet:   data["PreferencesSet"].(bool),
		LastUsed:         data["LastUsed"].(int64),
		FreeSentences:    int(data["FreeSentences"].(int64)),
	}, nil
}

// SetUserSentenceLanguage updates user's language of generated sentences
func (store *Store) SetUserSentenceLanguage(ctx context.Context, chatId int64, sentenceLanguage string) error {
	_, err := store.db.Collection("users").Doc(strconv.Itoa(int(chatId))).Update(ctx, []firestore.Update{
		{
			Path:  "SentenceLanguage",
			Value: sentenceLanguage,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserPremium updates user premiumUntil field to a new time stamp provided in unix time format
func (store *Store) UpdateUserPremium(ctx context.Context, chatId int64, premiumUntil int64) error {
	_, err := store.db.Collection("users").Doc(strconv.Itoa(int(chatId))).Update(ctx, []firestore.Update{
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

// SetUserLevel sets language level (e.g. A1, B2) for sentences that user will generate
// Also sets preferencesSet field to true because setting language is the last step of preferences
func (store *Store) SetUserLevel(ctx context.Context, chatId int64, level string) error {
	_, err := store.db.Collection("users").Doc(strconv.Itoa(int(chatId))).Update(ctx, []firestore.Update{
		{
			Path:  "Level",
			Value: level,
		},
		{
			Path:  "PreferencesSet",
			Value: true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
