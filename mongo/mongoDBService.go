package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"ibuYemekBotu/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB() *mongo.Collection {
	fmt.Println("Connecting to MongoDB")
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("Client Set Failed: ", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Connection to MongoDB Failed: ", err)
	}

	fmt.Println("Connected to MongoDB!")
	testDB := client.Database("ibuYemekBotu").Collection("users")

	return testDB
}

func Adduser(user *models.User) {
	collection := connectDB()
	insertResult, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted 1 document: ", insertResult.InsertedID)
}

func GetAllUsers() []models.User {
	collection := connectDB()
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Hata : " + err.Error())
	}
	defer cur.Close(ctx)
	var list []models.User
	for cur.Next(ctx) {
		var result models.User
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal("Hata : " + err.Error())
		}
		list = append(list, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal("Hata : " + err.Error())
	}
	fmt.Println(list)
	return list
}

func GetUser(chatid int64) bool {
	collection := connectDB()
	var result models.User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, bson.M{"chatid": chatid}).Decode(&result)
	if err != nil {
		return false
		//log.Fatal("Search Error : " + err.Error())
	}
	out, err := json.Marshal(&result)
	if err != nil {
		return false
		//log.Fatal("No Found : " + err.Error())

	}
	fmt.Println("Found User: " + string(out))
	return true
}

func DeleteUser(chatid int64) {
	collection := connectDB()
	_, err := collection.DeleteOne(context.TODO(), bson.D{{"chatid", chatid}})
	if err != nil {
		log.Fatal("Hata : " + err.Error())
	}
	fmt.Println("User deleted")
}
