package api_users

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kilfu0701/echo_app/core"
)

const TableName = "api_users"

type ApiUsers struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

type LastLogin struct {
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
	hour      string `bson:"hour" json:"hour"`
	Day       string `bson:"day" json:"day"`
	Month     string `bson:"month" json:"month"`
	Year      string `bson:"year" json:"year"`
}

type ApiUsersDoc struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Id2            string             `json:"id2"`
	UserName       string             `bson:"user_name" json:"user_name"`
	UserLevel      int64              `bson:"user_level" json:"user_level"`
	UserKey        string             `bson:"user_key" json:"user_key"`
	UserSecretHash string             `bson:"user_secret_hash" json:"user_secret_hash"`
	Mtime          int64              `bson:"mtime" json:"mtime"`
	Ctime          int64              `bson:"ctime" json:"ctime"`
	Status         int64              `bson:"status" json:"status"`
	LastLogin      *LastLogin         `bson:"last_login" json:"last_login"`
	//Owner string `bson:"owner,omitempty" json:"owner,omitempty"`
}

func New(ac core.AppContext) *ApiUsers {
	return &ApiUsers{
		Ctx:        ac.AppDBCtx,
		Collection: ac.AppDB.Collection(TableName),
	}
}

func (au *ApiUsers) FindAll() (*[]ApiUsersDoc, error) {
	filter := bson.M{"status": 0}
	var results []ApiUsersDoc

	cur, err := au.Collection.Find(au.Ctx, filter)
	if err != nil {
		log.Printf("err = %+v", err)
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var doc ApiUsersDoc
		if err := cur.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		doc.Id2 = core.EncodeIdToBase64(doc.Id.Hex())
		results = append(results, doc)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &results, nil
}

func (au *ApiUsers) FindById(id primitive.ObjectID) (*ApiUsersDoc, error) {
	filter := bson.M{
		"_id":    id,
		"status": 0,
	}
	var doc ApiUsersDoc

	err := au.Collection.FindOne(au.Ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (au *ApiUsers) FindByName(username string) (*ApiUsersDoc, error) {
	filter := bson.M{
		"user_name": username,
		"status":    0,
	}

	var doc ApiUsersDoc

	if err := au.Collection.FindOne(au.Ctx, filter).Decode(&doc); err != nil {
		log.Printf("err = %+v", err)
		return nil, err
	}

	doc.Id2 = core.EncodeIdToBase64(doc.Id.Hex())

	return &doc, nil
}

func (au *ApiUsers) FindByKey(key string) (*ApiUsersDoc, error) {
	filter := bson.M{
		"user_key": key,
		"status":   0,
	}

	var doc ApiUsersDoc

	if err := au.Collection.FindOne(au.Ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}

	doc.Id2 = core.EncodeIdToBase64(doc.Id.Hex())

	return &doc, nil
}

func (au *ApiUsers) CreateUser(username string, userSecretHash string) (*ApiUsersDoc, error) {
	// generate user_key
	uid := core.Uniqid("new_user", true)
	ts := core.Microtime()
	input_str := fmt.Sprintf("%s%s%s", ts, uid, username)
	h := sha1.New()
	io.WriteString(h, input_str)
	hashed := string(h.Sum(nil))
	hashedString := fmt.Sprintf("%x", hashed)
	userKey := core.UrlsafeBase64Encode(hashedString)

	now := time.Now()

	result, err := au.Collection.InsertOne(
		context.Background(),
		bson.D{
			{"user_name", username},
			{"user_level", 0},
			{"user_key", userKey},
			{"user_secret_hash", userSecretHash},
			{"owner", nil},
			{"mtime", now.Unix()},
			{"ctime", now.Unix()},
			{"status", 0},
		})

	if err != nil {
		return nil, err
	}

	if result.InsertedID == nil {
		return nil, errors.New("InsertedID is nil.")
	}

	log.Printf("InsertedID = %++v", result.InsertedID)
	doc, err := au.FindById(result.InsertedID.(primitive.ObjectID))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

/*
   "last_login": {
        "timestamp": NumberInt(1553735698),
        "hour": "2019-03-28 10",
        "day": "2019-03-28",
        "month": "2019-03",
        "year": "2019"
   }
*/
func (au *ApiUsers) UpdateLastLogin(id primitive.ObjectID) error {
	doc, err := au.FindById(id)
	if err != nil {
		return err
	}

	now := time.Now()

	result, err := au.Collection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": doc.Id,
		},
		bson.D{
			{
				"$set", bson.D{
					{
						"last_login", bson.D{
							{"timestamp", now.Unix()},
							{"hour", now.Format("2006-01-02 15")},
							{"day", now.Format("2006-01-02")},
							{"month", now.Format("2006-01")},
							{"year", now.Format("2006")},
						},
					},
				},
			},
		},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount != 1 && result.ModifiedCount != 1 {
		return errors.New("update failed")
	}

	return nil
}
