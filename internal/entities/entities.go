package entities

type MongoUser struct {
	Name    string `json:"name" bson:"name"`
	Age     int    `json:"age" bson:"age"`
	Friends []int  `json:"friends" bson:"friends"`
	ID      int    `json:"id" bson:"id"`
}
