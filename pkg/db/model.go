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
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	BID                    string             `bson:"bid"`
	CID                    string             `bson:"cid"`
	DispatchCenter         string             `bson:"dispatch_center"`
	DispatchCounterNo      string             `bson:"dispatch_center_counter_no"`
	ReceiptCenter          string             `bson:"receipt_center"`
	ReceiptCounterNo       string             `bson:"receipt_center_counter_no"`
	BusNo                  string             `bson:"bus_no"`
	DriverName             string             `bson:"driver_name"`
	DriverContact          string             `bson:"driver_contact"`
	LocationIncharge       string             `bson:"location_incharge"`
	LocationInchargeContact string            `bson:"location_incharge_contact"`
	FoodBeddingIncharge    string             `bson:"food_bedding_incharge"`
	FoodBeddingContact     string             `bson:"food_bedding_incharge_contact"`
	VolunteerName          string             `bson:"volunteer_name"`
	VolunteerContact       string             `bson:"volunteer_contact"`
	MicroObserverName      string             `bson:"micro_observer_name"`
	MicroObserverContact   string             `bson:"micro_observer_contact"`
	PollingStaff           []Staff            `bson:"polling_staff"`
	PoliceStaff            []Staff            `bson:"police_staff"`
}

type Staff struct {
	Name    string `bson:"name"`
	Contact string `bson:"contact"`
}
type Booth struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Cid               string             `bson:"cid"`
	Bid               string             `bson:"bid"`
	BoothName        string             `bson:"booth_name"`
	BloName           string             `bson:"blo_name"`
	BloContact       string             `bson:"blo_contact"`
	GeoCoordinates   string             `bson:"geo_coordinates"`
	SectorOfficerName string             `bson:"sector_officer_name"`
	SectorOfficerContact string             `bson:"sector_officer_contact"`
	AssistantSectorOfficerName string             `bson:"assistant_sector_officer_name"`
	AssistantSectorOfficerContact string             `bson:"assistant_sector_officer_contact"`
	LocationInchargeName  string             `bson:"location_incharge_name"`
	LocationInchargeContact string             `bson:"location_incharge_contact"`
	AroName             string             `bson:"aro_name"`
	AroContact          string             `bson:"aro_contact"`
	DspName             string             `bson:"dsp_name"`
	DspContact          string             `bson:"dsp_contact"`
	ShoName             string             `bson:"sho_name"`
	ShoContact          string             `bson:"sho_contact"`
	ChowkiInchargeName  string             `bson:"chownki_incharge_name"`
	ChowkiInchargeContact string             `bson:"chownki_incharge_contact"`
	MedicalOfficerName  string             `bson:"medical_officer_name"`
	MedicalOfficerContact string             `bson:"medical_officer_contact"`
	SmoName             string             `bson:"smo_name"`
	SmoContact          string             `bson:"smo_contact"`
	AmbulanceContact    string             `bson:"ambulance_contact"`
	FireContact         string             `bson:"fire_contact"`
	HeatwaveName        string             `bson:"heatwave_name"`
	HeatwaveContact    string             `bson:"heatwave_contact"`
	Counter             string             `bson:"counter"`
	LastUpdated        string             `bson:"last_updated"`
  }
  
  type Voter struct {
	ID      string `bson:"_id,omitempty"`
	CID     string `bson:"cid"`
	BID     string `bson:"bid"`
	Name    string `bson:"name"`
	Contact string `bson:"contact"`
	Message string `bson:"message"`
	Status  string `bson:"status"`
}

type Counter struct {
	Counter string `bson:"counter"`
}