package indexing

import (
	"context"
	"fmt"
	"log"
	"studenent_manager/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Index_Student_ID            = "Student_ID"
	Index_Student_Name          = "Student_Name"
	Index_Student_Age           = "Index_Student_Age"
	Index_Student_Age_Name      = "Index_Student_Age_Name"
	Index_WildCard_Student      = "Index_WildCard_Student"
	Index_Partial_Student_Score = "Index_Partial_Student_Score"
)

func NewIndexingByStudentID() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "student_id", Value: 1}}, // 1 là ascending, -1 là descending
		Options: options.Index().SetUnique(true).SetName(Index_Student_ID),
	}
	return indexModel
}

func NewIndexingByStudentName() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName(Index_Student_Name),
	}

	return indexModel
}

func NewIndexingByStudentAge() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "age", Value: 1}},
		//Options: options.Index().SetUnique(true),
		Options: options.Index().SetName(Index_Student_Age),
	}

	return indexModel
}

func NewIndexingByStudentAgeAndName() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "age", Value: 1},
		},
		Options: options.Index().SetName(Index_Student_Age_Name),
		//Options: options.Index().SetUnique(true),
	}
	return indexModel
	// _, err := mgm.Coll(&db.Student{}).Indexes().CreateMany(context.Background(), []mongo.IndexModel{indexModel})
}

// Wildcard Index
func NewWildCardIndex() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "$**", Value: 1}}, // wildcard index
		Options: options.Index().SetName(Index_WildCard_Student),
	}
	// _, err := mgm.Coll(&db.Student{}).Indexes().CreateOne(context.Background(), indexModel)

	return indexModel
}

// Partial Index
func NewPartialIndexScore(targetScore int) mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "score", Value: 1}},
		Options: options.Index().SetPartialFilterExpression(bson.M{
			"score": bson.M{"$gt": targetScore},
		}).SetName(Index_Partial_Student_Score),
	}
	// _, err := mgm.Coll(&db.Student{}).Indexes().CreateOne(context.Background(), indexModel)
	return indexModel
}

func AddIndexModels(indexModels []mongo.IndexModel) error {
	_, err := mgm.Coll(&db.Student{}).Indexes().CreateMany(context.Background(), indexModels)
	return err
}

func AddIndexModel(indexModel mongo.IndexModel) error {
	_, err := mgm.Coll(&db.Student{}).Indexes().CreateOne(context.Background(), indexModel)
	return err
}

// 2D index --> use location (longtitude, latitude)
// func Create2DIndex() {
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{{Key: "location", Value: "2d"}}, // 2D index
// 	}
// 	indexName, err := studentCollection.Indexes().CreateOne(context.TODO(), indexModel)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Created 2D Index:", indexName)
// }

// 2D sphere index --> use in geoJSON
// func Create2DSphereIndex() {
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{{Key: "geoLocation", Value: "2dsphere"}},
// 	}
// 	indexName, err := studentCollection.Indexes().CreateOne(context.TODO(), indexModel)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Created 2D Sphere Index:", indexName)
// }

// B+tree

func DeleteAllIndexing() error {
	_, err := mgm.Coll(&db.Student{}).Indexes().DropAll(context.TODO())
	return err
}

func DeleteIndexing(indexName string) error {
	_, err := mgm.Coll(&db.Student{}).Indexes().DropOne(context.TODO(), indexName)
	return err
}

func PrintAllIndexing() {
	cursor, err := mgm.Coll(&db.Student{}).Indexes().List(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var indexes []bson.M
	if err := cursor.All(context.TODO(), &indexes); err != nil {
		log.Fatal(err)
	}
	for _, index := range indexes {
		fmt.Println(index)
	}
}
