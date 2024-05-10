package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
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

func AuthenticateAdmin(username, password, usertype string) (string, error) {
	var userAro ARO
	var userBlo BLO
	if usertype == "aro" {
		filter := bson.M{"user": usertype, "aro_name": username, "pass": password}
		err := usersCollection.FindOne(context.Background(), filter).Decode(&userAro)
		if err != nil {
			return "", err
		}
		return userAro.CID, nil
	} else if usertype == "blo" {
		filter := bson.M{"user": usertype, "blo_name": username, "pass": password}
		err := usersCollection.FindOne(context.Background(), filter).Decode(&userBlo)
		if err != nil {
			return "", err
		}
		return userBlo.BID, nil
	} else {
		return "", nil
	}
}

func GetLiveTrafficByBoothID(boothID string) (int, error) {
	var liveTraffic = 0

	filter := bson.M{"bid": boothID}
	result := pollingStationCollection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return liveTraffic, fmt.Errorf("no live traffic data found for the specific BOOTH ID : %v", boothID)
		}
		return liveTraffic, fmt.Errorf("error finding live traffic data %v", result.Err())
	}

	if err := result.Decode(&liveTraffic); err != nil {
		return liveTraffic, fmt.Errorf("error decoding live traffic data %v", err)
	}

	return liveTraffic, nil
}
