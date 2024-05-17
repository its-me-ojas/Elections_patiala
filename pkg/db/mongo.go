package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pollingStationCollection *mongo.Collection
var usersCollection *mongo.Collection
var votersReqCollection *mongo.Collection
var displayDataConnection *mongo.Collection

func InitDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	dbURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Mongodb")
	pollingStationCollection = client.Database("election_patiala").Collection("polling_station")
	usersCollection = client.Database("election_patiala").Collection("users")
	votersReqCollection = client.Database("election_patiala").Collection("voters_req")
	displayDataConnection = client.Database("election_patiala").Collection("display_data")
}

func AuthenticateAdmin(contact, password string) (map[string]string, error) {
	var user User
	filter := bson.M{"contact": contact, "pass": password}
	err := usersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
        return nil, err
    }
	result := make(map[string]string)
    result["usertype"] = user.UserType

    switch user.UserType {
    case "aro":
        result["cid"] = user.CID
    case "blo", "ps":
        result["cid"] = user.CID
        result["bid"] = user.BID
    default:
        return nil, nil
    }

    return result, nil
}

func GetLiveTrafficByBoothID(boothID string) (int, error) {
	var liveTraffic struct {
		Counter int `bson:"counter"`
	}

	filter := bson.M{"bid": boothID}
	result := pollingStationCollection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return liveTraffic.Counter, fmt.Errorf("no live traffic data found for the specific BOOTH ID : %v", boothID)
		}
		return liveTraffic.Counter, fmt.Errorf("error finding live traffic data %v", result.Err())
	}

	if err := result.Decode(&liveTraffic); err != nil {
		return liveTraffic.Counter, fmt.Errorf("error decoding live traffic data %v", err)
	}

	return liveTraffic.Counter, nil
}
