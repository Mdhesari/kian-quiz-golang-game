package main

import (
	"context"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	// "os/user"

	// "mdhesari/kian-quiz-golang-game/entity"

	// "mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	url := "mongodb://michael:secret@localhost:27017/"
	clientOptions := options.Client().ApplyURI(url)

	cli, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	if err = cli.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	collection := cli.Database("freedom").Collection("users")

	user := entity.User{
		Name:     "Yalda",
		Email:    "yalda@aqafazli.ir",
		Mobile:   "09222222222",
		Password: []byte("mahsaamini"),
	}

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println("Happy Birthday")

	cur, err := collection.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		panic(err)
	}

	for cur.Next(context.Background()) {
		var u entity.User

		if err := cur.Decode(&u); err != nil {
			panic(err)
		}

		fmt.Println(u)
	}

	log.Println("MongoClient Connected...")
}
