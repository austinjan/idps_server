package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB struct of monge database
type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

var instance *DB
var once sync.Once

const defautlURL string = "mongodb://localhost:27017"

// GetDB get unique database pointer
func GetDB() *DB {
	once.Do(func() {
		instance = initDB(defautlURL)
	})
	return instance
}

// Connect to default database
func (d *DB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := d.client.Connect(ctx)
	if err != nil {
		return err
	}
	err = d.client.Ping(context.TODO(), nil)
	return err
}

// Initial databas
func initDB(url string) *DB {
	const defaultDatabase = "test"
	db := new(DB)
	var err error
	db.client, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = db.client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = db.client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	db.db = db.client.Database(defaultDatabase)
	fmt.Printf("Database %s connected... \n", url)
	return db
}

func (d *DB) SaveTagPosition(data bson.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection("tagPosition")

	_, err := collection.InsertOne(ctx, data)
	//fmt.Println("Insert id:", res.InsertedID)

	if err != nil {
		log.Println(err)
	}

}

// Insert insert data into collection c
func (d *DB) Insert(c string, data bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)

	res, err := collection.InsertOne(ctx, data)
	fmt.Println("Insert id:", res.InsertedID)

	if err != nil {
		log.Println(err)
	}
	return err

}

func (d *DB) SaveTagInfo(data bson.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection("tagInformation")

	_, err := collection.InsertOne(ctx, data)
	//fmt.Println("Insert id:", res.InsertedID)

	if err != nil {
		log.Println(err)
	}

}

func (d *DB) read(c string, query bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	cur, err := collection.Find(ctx, query)

	var data []bson.M
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("collection: %s is not exist!", c)
		return nil, errors.New(msg)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var item bson.M
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	if len(data) == 0 {
		msg := fmt.Sprintf("Query: %s not found result.", c)
		return nil, errors.New(msg)
	}
	return data, nil
}

// Read - db.c.find(query) -> return data "ARRAY" of json format , error
func (d *DB) Read(c string, query bson.M) ([]byte, error) {
	fmt.Println("Read query", query)
	data, err := d.read(c, query)
	if err != nil {
		fmt.Println("Read err", err)
		return nil, err
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Read json err", err)
		return nil, err
	}
	return jsonStr, err
}
