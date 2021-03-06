package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type RssFields struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"pubDate"`
	CreatedAt   time.Time `json:"createdDate"`
}

type TweetFields struct {
	Name        string    `json:"name"`
	User        string    `json:"user"`
	Id          string    `json:"id"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"pubDate"`
	CreatedAt   time.Time `json:"createdDate"`
}

type PublishFields struct {
	Description string    `json:"description"`
	ImagePath   string    `json:"imagepath"`
	CreatedAt   time.Time `json:"createdDate"`
}

type Storage struct {
	c     mongo.Client
	dburl string
	db    *mongo.Database
	dbid  string
}

var dbUrl string
var dbId string
var dbUser string
var dbPasswd string
var sharedInstance *Storage

func newStorage() *Storage {
	if dbUrl == "" {
		return nil
	}

	fmt.Println("newStorage: ", dbUrl, dbId)
	c := connect(dbUrl, dbId, dbUser, dbPasswd)
	db := c.Database(dbId)
	return &Storage{c, dbUrl, db, dbId}
}

func SetConfig(dburl string, dbid string, dbuser string, dbpasswd string) {
	dbUrl = dburl
	dbId = dbid
	dbUser = dbuser
	dbPasswd = dbpasswd
}

func GetInstance() *Storage {
	if sharedInstance == nil {
		sharedInstance = newStorage()
	}
	return sharedInstance
}

func TermInstance() {
	if sharedInstance != nil {
		disconnect(sharedInstance.c)
	}
}

func connect(dburl string, dbid string, dbuser string, dbpasswd string) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dburl).
		SetAuth(options.Credential{AuthSource: dbid, Username: dbuser, Password: dbpasswd})
	c, err := mongo.Connect(ctx, clientOptions)
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("connection error:", err)
	} else {
		fmt.Println("connection success:")
	}
	return *c
}

func disconnect(c mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.Disconnect(ctx)
	if err != nil {
		fmt.Println("disconnection error:", err)
	} else {
		fmt.Println("disconnection success:")
	}
}

func InsertRSS(s *Storage, colid string, title string, link string, desc string, pubdate time.Time) error {
	col := s.db.Collection(colid)

	doc := RssFields{
		Title:       title,
		Link:        link,
		Description: desc,
		PublishedAt: pubdate,
		CreatedAt:   time.Now(),
	}

	_, err := col.InsertOne(context.Background(), doc)
	return err
}

func FindRSS(s *Storage, colid string, link string, pubdate time.Time) (bool, error) {
	col := s.db.Collection(colid)

	filter := struct {
		Link        string
		PublishedAt time.Time
	}{link, pubdate}

	var doc bson.Raw
	findOptions := options.FindOne()
	err := col.FindOne(context.Background(), filter, findOptions).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, err
	}

	return true, err
}

func InsertTweet(s *Storage, colid string, name string, user string, id string, link string, desc string, pubdate time.Time) error {
	col := s.db.Collection(colid)

	doc := TweetFields{
		Name:        name,
		User:        user,
		Id:          id,
		Link:        link,
		Description: desc,
		PublishedAt: pubdate,
		CreatedAt:   time.Now(),
	}

	_, err := col.InsertOne(context.Background(), doc)
	return err
}

func FindTweet(s *Storage, colid string, user string, id string) (bool, error) {
	col := s.db.Collection(colid)

	filter := struct {
		User string
		Id   string
	}{user, id}

	var doc bson.Raw
	findOptions := options.FindOne()
	err := col.FindOne(context.Background(), filter, findOptions).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, err
	}

	return true, err
}

func InsertPublish(s *Storage, colid string, desc string, imgpath string) error {
	col := s.db.Collection(colid)

	doc := PublishFields{
		Description: desc,
		ImagePath:   imgpath,
		CreatedAt:   time.Now(),
	}

	_, err := col.InsertOne(context.Background(), doc)
	return err
}

func FindPublish(s *Storage, colid string) (string, string, error) {
	col := s.db.Collection(colid)

	var doc struct {
		Description string             `json:"description"`
		ImagePath   string             `json:"imagepath"`
		ID          primitive.ObjectID `json:"id" bson:"_id"`
	}
	findOptions := options.FindOne()
	err := col.FindOne(context.Background(), bson.D{}, findOptions).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return "", "", err
	}
	desc := doc.Description
	imgpath := doc.ImagePath

	deleteOptions := options.Delete()
	_, err = col.DeleteOne(context.Background(), bson.M{"_id": doc.ID}, deleteOptions)
	if err != nil {
		fmt.Println("delete one error:", err)
		return "", "", err
	}

	return desc, imgpath, err
}
