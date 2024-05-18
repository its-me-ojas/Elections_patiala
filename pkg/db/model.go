package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"id,omitempty"`
	UserType  string             `bson:"user"`
	CID       string             `bson:"cid,omitempty"`
	BID       string             `bson:"bid,omitempty"`
	Name      string             `bson:"blo_name"`
	BoothName string             `bson:"booth_name"`
	Contact   string             `bson:"contact"`
	Password  string             `bson:"pass"`
}

type DisplayData struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	BID                   string             `bson:"bid"`
	CID                   string             `bson:"cid"`
	DispatchCenter        string             `bson:"dispatch_center"`
	CounterNoBase         string             `bson:"counter_no_base"`
	BusNo                 string             `bson:"bus_no"`
	DriverNo              string             `bson:"driver_no"`
	DriverName            string             `bson:"driver_name"`
	SchoolIncharge        string             `bson:"school_incharge"`
	SchoolInchargeNumber  string             `bson:"school_incharge_number"`
	HousingIncharge       string             `bson:"housing_incharge"`
	HousingInchargeNumber string             `bson:"housing_incharge_number"`
	PollingStaff          []Staff            `bson:"polling_staff"`
	PoliceStaff           []Staff            `bson:"police_staff"`
	Microbserver          string             `bson:"microbserver"`
	MicrobserverContact   string             `bson:"microbserver_contact"`
}

type Staff struct {
	Name    string `bson:"name"`
	Contact string `bson:"contact"`
}
