package main

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
	CreatedAt   time.Time `json:"createdDate"`
}

func DBConnect(dburl string) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("connection error:", err)
	} else {
		fmt.Println("connection success:")
	}
	return *c
}

func DBDisconnect(c mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.Disconnect(ctx)
	if err != nil {
		fmt.Println("disconnection error:", err)
	} else {
		fmt.Println("disconnection success:")
	}
}

func DBInsertRSS(col *mongo.Collection, title string, link string, desc string, pubdate time.Time) error {
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

func DBFindRSS(col *mongo.Collection, link string) (bool, error) {
	filter := struct {
		Link string
	}{link}

	var doc bson.Raw
	findOptions := options.FindOne()
	err := col.FindOne(context.Background(), filter, findOptions).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, err
	}

	return true, err
}

func DBInsertTweet(col *mongo.Collection, name string, user string, id string, link string, desc string, pubdate time.Time) error {
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

func DBFindTweet(col *mongo.Collection, user string, id string) (bool, error) {
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

func DBInsertPublish(col *mongo.Collection, desc string) error {
	doc := PublishFields{
		Description: desc,
		CreatedAt:   time.Now(),
	}

	_, err := col.InsertOne(context.Background(), doc)
	return err
}

func DBFindPublish(col *mongo.Collection) (string, error) {
	var doc struct {
		Description string             `json:"description"`
		ID          primitive.ObjectID `json:"id" bson:"_id"`
	}
	findOptions := options.FindOne()
	err := col.FindOne(context.Background(), bson.D{}, findOptions).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return "", err
	}
	desc := doc.Description

	deleteOptions := options.Delete()
	_, err = col.DeleteOne(context.Background(), bson.M{"_id": doc.ID}, deleteOptions)
	if err != nil {
		fmt.Println("delete one error:", err)
		return "", err
	}

	return desc, err
}
