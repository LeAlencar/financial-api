package models

type PayloadResponse struct {
	Code       string `json:"code" bson:"code"`
	Codein     string `json:"codein" bson:"codein"`
	Name       string `json:"name" bson:"name"`
	High       string `json:"high" bson:"high"`
	Low        string `json:"low" bson:"low"`
	VarBid     string `json:"varBid" bson:"varBid"`
	PctChange  string `json:"pctChange" bson:"pctChange"`
	Bid        string `json:"bid" bson:"bid"`
	Ask        string `json:"ask" bson:"ask"`
	Timestamp  string `json:"timestamp" bson:"timestamp"`
	CreateDate string `json:"create_date" bson:"create_date"`
}
