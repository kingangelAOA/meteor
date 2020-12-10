package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"meteor/configs"
	"testing"
	"time"
)

func Test(t *testing.T) {

}

func TestMongo(t *testing.T) {
	newC := configs.Mongo{
		Host:     "127.0.0.1",
		Port:     "27017",
		Db:       "datareport",
	}
	oldC := configs.Mongo{
		Host: "10.40.181.58",
		Port: "20000",
		Db:   "datareport",
	}
	clientOld := CreateMongo(oldC)
	clientNew := CreateMongo(newC)
	ctx, cancel := context.WithTimeout(context.Background(), 1000000*time.Second)
	defer cancel()
	if err := clientOld.Client.Connect(ctx); err != nil {
		fmt.Println(err)
	}
	if err := clientNew.Client.Connect(ctx); err != nil {
		fmt.Println(err)
	}
	clientNew.Client.Database("datareport")
	implementsOld := clientOld.Client.Database("datareport").Collection("algorithm_cases")
	implementsNew := clientNew.Client.Database("datareport").Collection("algorithm_cases")
	fmt.Println(implementsNew)
	cur, err := implementsOld.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}
	defer cur.Close(ctx)
	var results []interface{}
	count := 0
	for cur.Next(ctx) {
		var result bson.M

		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
		// delete(result, "_id")
		// fmt.Println(result)

		// do something with result....
		if count%50 == 0 && count > 0 {
			startT := time.Now()
			if _, err := implementsNew.InsertMany(ctx, results); err != nil {
				fmt.Println(err)
				panic(err)
			}
			tc := time.Since(startT)
			fmt.Printf("time cost = %v\n: %v", tc, count)
			fmt.Println("**************", count)
			results = nil
		}

		count++
	}

	if len(results) > 0 {
		if r, err := implementsNew.InsertMany(ctx, results); err != nil {
			fmt.Println(err)
			panic(err)
		} else {
			fmt.Println(r)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}

func TestMongoA(t *testing.T) {
	//oldC := configs.Mongo{
	//	User:     "helix_dashboard",
	//	Password: "123456",
	//	Host:     "10.40.181.235",
	//	Port:     "27017",
	//	Db:       "dashboard",
	//}
	//clientOld := MongoInit(oldC)
	//ctx, cancel := context.WithTimeout(context.Background(), 1000000*time.Second)
	//defer cancel()
	//if err := clientOld.Connect(ctx); err != nil {
	//	fmt.Println(err)
	//}
	//
	//implementsOld := clientOld.Database("dashboard").Collection("TXWL_Bug")
	//opts := options.Aggregate().SetMaxTime(15 * time.Second)
	//matchStage := bson.D{{"$match", bson.D{{"status.name", bson.D{{"$ne", "Closed"}}}}}}
	//s := time.Now()
	//cur, err := implementsOld.Aggregate(context.TODO(), mongo.Pipeline{matchStage}, opts)
	//fmt.Println(time.Since(s))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer cur.Close(ctx)
	//count := 0
	//for cur.Next(ctx) {
	//
	//	count++
	//}
	//
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
}
