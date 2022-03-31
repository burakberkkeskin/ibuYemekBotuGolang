package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"ibuYemekBotu/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	getEnv()
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
	//testDB := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))

	return client
}

var Client *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

var userCollection *mongo.Collection = GetCollection(Client, os.Getenv("DATABASE"), os.Getenv("COLLECTION"))

func Adduser(user *models.User) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error While Inserting User: ", err)
	}

	fmt.Println("Inserted 1 User: ", insertResult.InsertedID)
}

func GetAllUsers() *[]models.User {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := userCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Error While Getting All Users: " + err.Error())
	}

	defer cur.Close(ctx)

	var userList []models.User
	for cur.Next(ctx) {
		var result models.User
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal("Error While Decoding Results While Getting All Users : " + err.Error())
		}
		userList = append(userList, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal("Error While Getting All Users : " + err.Error())
	}

	return &userList
}

func GetUser(chatid int64) bool {

	var result models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"chatid": chatid}).Decode(&result)
	if err != nil {
		log.Println("Search Error : " + err.Error())
		return false
	}

	out, err := json.Marshal(&result)
	if err != nil {
		log.Println("No User Found : " + err.Error())
		return false

	}
	fmt.Println("User Found: " + string(out))
	return true
}

func DeleteUser(chatid int64) bool {
	_, err := userCollection.DeleteOne(context.TODO(), bson.D{{Key: "chatid", Value: chatid}})
	if err != nil {
		log.Println("Hata : " + err.Error())
		return false
	}
	fmt.Println("User deleted")
	return true
}

func getEnv() {
	err := godotenv.Load("./configs/.env.dev")
	if err != nil {
		log.Println("Error loading .env file")
	}
}
