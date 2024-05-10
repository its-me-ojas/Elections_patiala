package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type ARO struct {
	ID         primitive.ObjectID `bson:"id,omitempty"`
	CID        string             `bson:"cid,omitempty"`
	AROName    string             `bson:"aro_name"`
	AROContact string             `bson:"aro_contact"`
	Password   string             `bson:"password"`
}

type BLO struct {
	ID         primitive.ObjectID `bson:"id,omitempty"`
	CID        string             `bson:"cid,omitempty"`
	BID        string             `bson:"bid,omitempty"`
	BLOName    string             `bson:"blo_name"`
	BoothName  string             `bson:"booth_name"`
	BLOContact string             `bson:"blo_contact"`
	Password   string             `bson:"password"`
}
