package tests

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/HiChen85/RedditRecipesWithGin/recipes_service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

var (
	MONGO_USERNAME   = "admin"
	MONGO_PASSWORD   = "password"
	MONGO_HOST       = "localhost"
	MONGO_PORT       = "27017"
	MONGO_DATABASE   = "demo"
	MONGO_COLLECTION = "users"
)

func TestAddUser(t *testing.T) {
	users := map[string]string{
		"admin":      "fCRmh4Q2J7Rseqkz",
		"packt":      "RE4zfHB35VPtTkbT",
		"mlabouardy": "L3nSFRcZzNQ67bcc",
	}
	var MONGO_URI = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=%v", MONGO_USERNAME, MONGO_PASSWORD, MONGO_HOST, MONGO_PORT, MONGO_DATABASE, MONGO_USERNAME)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		t.Fatal(err)
	}
	if err = mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		t.Fatal(err)
	}
	collection := mongoClient.Database(MONGO_DATABASE).Collection(MONGO_COLLECTION)
	h := sha256.New()
	for user, pwd := range users {
		one, err := collection.InsertOne(context.TODO(), bson.M{
			"username": user,
			"password": string(h.Sum([]byte(pwd))),
		})
		if err != nil {
			return
		}
		t.Log(one.InsertedID)
	}
}

func TestGetUser(t *testing.T) {
	users := map[string]string{
		"admin":      "fCRmh4Q2J7Rseqkz",
		"packt":      "RE4zfHB35VPtTkbT",
		"mlabouardy": "L3nSFRcZzNQ67bcc",
	}
	var MONGO_URI = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=%v", MONGO_USERNAME, MONGO_PASSWORD, MONGO_HOST, MONGO_PORT, MONGO_DATABASE, MONGO_USERNAME)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		t.Fatal(err)
	}
	if err = mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		t.Fatal(err)
	}
	collection := mongoClient.Database(MONGO_DATABASE).Collection(MONGO_COLLECTION)
	h := sha256.New()
	for u, p := range users {
		var user models.User
		cur := collection.FindOne(context.TODO(), bson.M{
			"username": u,
			"password": string(h.Sum([]byte(p))),
		})
		if cur.Err() != nil {
			t.Log("没查到", u)
		}
		cur.Decode(&user)
		t.Log(user)
	}
}
