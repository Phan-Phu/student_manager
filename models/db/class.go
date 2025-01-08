package db

type Class struct {
	ClassID    int       `json:"id" bson:"id"`
	ClassName  int       `json:"name" bson:"name"`
	Students   []Student `bson:"-"` // not need store in database
	Teachers   []Teacher `bson:"-"` // not need store in database
	CreateDate string    `json:"create_date" bson:"create_date"`
}
