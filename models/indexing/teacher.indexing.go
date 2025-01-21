package indexing

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Index_Teacher_ID            = "Teacher_ID"
	Index_Teacher_Name          = "Teacher_Name"
	Index_Teacher_Age           = "Index_Teacher_Age"
	Index_Teacher_Age_Name      = "Index_Teacher_Age_Name"
	Index_WildCard_Teacher      = "Index_WildCard_Teacher"
	Index_Partial_Teacher_Score = "Index_Partial_Teacher_Score"
)

func NewIndexingTeacher() {
	indexModels := []mongo.IndexModel{}
	// DeleteIndexing("name_index")
	DeleteAllStudentIndexing()

	indexModels = append(indexModels, NewIndexingByTeacherName())
	indexModels = append(indexModels, NewIndexingByTeacherAge())
	indexModels = append(indexModels, NewPartialIndexScore(50)) // target score: 50
	indexModels = append(indexModels, NewIndexingByTeacherAgeAndName())
	indexModels = append(indexModels, NewWildCardIndex())
	AddIndexModels(indexModels)
	PrintAllStudentIndexing()
}

func NewIndexingByTeacherName() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName(Index_Teacher_Name),
	}

	return indexModel
}

func NewIndexingByTeacherAge() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "age", Value: 1}},
		//Options: options.Index().SetUnique(true),
		Options: options.Index().SetName(Index_Teacher_Age),
	}

	return indexModel
}

func NewIndexingByTeacherAgeAndName() mongo.IndexModel {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "age", Value: 1},
		},
		Options: options.Index().SetName(Index_Teacher_Age_Name),
		//Options: options.Index().SetUnique(true),
	}
	return indexModel
	// _, err := mgm.Coll(&db.Teacher{}).Indexes().CreateMany(context.Background(), []mongo.IndexModel{indexModel})
}

// func AddIndexModels(indexModels []mongo.IndexModel) error {
// 	_, err := mgm.Coll(&db.Teacher{}).Indexes().CreateMany(context.Background(), indexModels)
// 	return err
// }

// func AddIndexModel(indexModel mongo.IndexModel) error {
// 	_, err := mgm.Coll(&db.Teacher{}).Indexes().CreateOne(context.Background(), indexModel)
// 	return err
// }

// 2D index --> use location (longtitude, latitude)
// func Create2DIndex() {
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{{Key: "location", Value: "2d"}}, // 2D index
// 	}
// 	indexName, err := TeacherCollection.Indexes().CreateOne(context.TODO(), indexModel)
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
// 	indexName, err := TeacherCollection.Indexes().CreateOne(context.TODO(), indexModel)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Created 2D Sphere Index:", indexName)
// }

// B+tree

// func DeleteAllIndexing() error {
// 	_, err := mgm.Coll(&db.Teacher{}).Indexes().DropAll(context.TODO())
// 	return err
// }

// func DeleteIndexing(indexName string) error {
// 	_, err := mgm.Coll(&db.Teacher{}).Indexes().DropOne(context.TODO(), indexName)
// 	return err
// }

// func PrintAllIndexing() {
// 	cursor, err := mgm.Coll(&db.Teacher{}).Indexes().List(context.TODO())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var indexes []bson.M
// 	if err := cursor.All(context.TODO(), &indexes); err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, index := range indexes {
// 		fmt.Println(index)
// 	}
// }
