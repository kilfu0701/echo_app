package user

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"app/core"
)

const TableName = "users"

type AppUser struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

type AppUserDoc struct {
	Id  primitive.ObjectID `bson:"_id" json:"_id"`
	Id2 string             `json:"id2"`
}

func New(ac core.AppContext) *AppUser {
	return &AppUser{
		Ctx:        ac.AppDBCtx,
		Collection: ac.AppDB.Collection(TableName),
	}
}

func (au *AppUser) FindAll() (*[]AppUserDoc, error) {
	filter := bson.M{"status": 0}
	var results []AppUserDoc

	cur, err := au.Collection.Find(au.Ctx, filter)
	if err != nil {
		log.Printf("err = %+v", err)
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var doc AppUserDoc
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

/*
func (au *AppUser) FindByName(username string) (*AppUserDoc, error) {
	filter := bson.M{

		"status": 0,
	}
}
*/
