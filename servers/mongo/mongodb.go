package mongodb

import (
	"context"
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
	collection := d.db.Collection("tagInformation")

	_, err := collection.InsertOne(ctx, data)
	//fmt.Println("Insert id:", res.InsertedID)

	if err != nil {
		log.Println(err)
	}

}
