package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"

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
	case "blo", "ps","vl":
		result["cid"] = user.CID
		result["bid"] = user.BID
	default:
		return nil, nil
	}

	return result, nil
}

func GetLiveTrafficByBoothID(boothID string) (string, error) {
	var liveTraffic struct {
		Counter string `bson:"counter"`
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

func UpdateQueue(cid, bid, counter string) error {
	filter := bson.M{"bid": bid, "cid": cid}
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Error loading IST location:", err)
		return err
	}
	currentTimeIST := time.Now().In(istLocation)

	update := bson.M{"$set": bson.M{"counter": counter, "last_updated": currentTimeIST.Format("1504")}}
	_, err = pollingStationCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePassword(contact, newPassword string) error {
	filter := bson.M{"contact": contact}
	
	update := bson.M{"$set": bson.M{"pass": newPassword}}
	_, err := usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetAllVoters(cid string) ([]Voter, error) {
	var voters []Voter

	filter := bson.M{"cid": cid}
	cursor, err := votersReqCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var voter Voter
		err := cursor.Decode(&voter)
		if err != nil {
			return nil, err
		}
		voters = append(voters, voter)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return voters, nil
}

func GetAllVotersBid(cid string, bid string) ([]Voter, error) {
	var voters []Voter

	filter := bson.M{"cid": cid, "bid": bid}
	cursor, err := votersReqCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var voter Voter
		err := cursor.Decode(&voter)
		if err != nil {
			return nil, err
		}
		voters = append(voters, voter)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return voters, nil
}


func GetQueue(cid, bid string) (string, error) {
	var counter Counter

	filter := bson.M{"bid": bid, "cid": cid}
	err := pollingStationCollection.FindOne(context.Background(), filter).Decode(&counter)
	if err != nil {
		return "", err
	}

	return counter.Counter, nil
}

func GetBooth(cid, bid string) (Booth, error) {
	var booth Booth
	filter := bson.M{"cid": cid, "bid": bid}
	err := pollingStationCollection.FindOne(context.Background(), filter).Decode(&booth)
	if err != nil {
		fmt.Println(err)
		return booth, err
	}

	return booth, nil
}
func GetDisplayData(cid, bid string) (DisplayData, error) {
	var display_data DisplayData
	filter := bson.M{"cid": cid, "bid": bid}
	err := displayDataConnection.FindOne(context.Background(), filter).Decode(&display_data)
	if err != nil {
		fmt.Println(err)
		return display_data, err
	}

	return display_data, nil
}

func UpdateVoterRequest(objectID string) error {
	objID, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		log.Println("Error converting string to ObjectID", err)
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": "resolved"}}
	_, err = votersReqCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating voter request", err)
		return err
	}
	return nil
}
func FetchBoothsByCidAndTime(cid string) ([]Booth, error) {
	
	filter := bson.M{"cid": cid}
	cur, err := pollingStationCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var booths []Booth
	for cur.Next(context.Background()) {
		var booth Booth
		if err := cur.Decode(&booth); err != nil {
			return nil, err
		}
		lastUpdated, err := time.Parse("1504", booth.LastUpdated)
		if err != nil {
			return nil, err
		}
		loc, err := time.LoadLocation("Asia/Kolkata")
        if err != nil {
            return nil, err
        }
        now := time.Now().In(loc)

        lastUpdated = time.Date(now.Year(), now.Month(), now.Day(), lastUpdated.Hour(), lastUpdated.Minute(), 0, 0, loc)

		if now.Sub(lastUpdated) > 45*time.Minute {
            booths = append(booths, booth)
        }
		
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return booths, nil
}
func GetAllStaffBid(cid string, bid string) ([]User, error) {
	var users []User

	filter := bson.M{"cid": cid, "bid": bid}
	cursor, err := usersCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}