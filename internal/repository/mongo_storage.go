package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

type MongoStorage struct {
	db *mongo.Database
}

func NewMongoStorage(uri, dbName string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoStorage{db: db}, nil
}

// Paste

func (m *MongoStorage) SavePaste(p model.Paste) error {
	_, err := m.db.Collection("pastes").InsertOne(context.TODO(), p)
	return err
}

func (m *MongoStorage) GetPasteByID(id string) (*model.Paste, error) {
	var paste model.Paste
	err := m.db.Collection("pastes").FindOne(context.TODO(), bson.M{"id": id}).Decode(&paste)
	if err != nil {
		return nil, err
	}
	return &paste, nil
}

func (m *MongoStorage) UpdatePaste(p model.Paste) error {
	_, err := m.db.Collection("pastes").ReplaceOne(context.TODO(), bson.M{"id": p.ID}, p)
	return err
}

func (m *MongoStorage) DeletePaste(id string) error {
	_, err := m.db.Collection("pastes").DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func (m *MongoStorage) GetAllPastes() ([]model.Paste, error) {
	cursor, err := m.db.Collection("pastes").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var pastes []model.Paste
	for cursor.Next(context.TODO()) {
		var p model.Paste
		if err := cursor.Decode(&p); err != nil {
			log.Println("decode error:", err)
			continue
		}
		pastes = append(pastes, p)
	}
	return pastes, nil
}

// User

func (m *MongoStorage) SaveUser(u model.User) error {
	_, err := m.db.Collection("users").InsertOne(context.TODO(), u)
	return err
}

func (m *MongoStorage) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := m.db.Collection("users").FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoStorage) UpdateUser(u model.User) error {
	_, err := m.db.Collection("users").ReplaceOne(context.TODO(), bson.M{"id": u.ID}, u)
	return err
}

func (m *MongoStorage) DeleteUser(id string) error {
	_, err := m.db.Collection("users").DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func (m *MongoStorage) GetAllUsers() ([]model.User, error) {
	cursor, err := m.db.Collection("users").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []model.User
	for cursor.Next(context.TODO()) {
		var u model.User
		if err := cursor.Decode(&u); err != nil {
			log.Println("decode error:", err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// ShortURL

func (m *MongoStorage) SaveShortURL(s model.ShortURL) error {
	_, err := m.db.Collection("shorturls").InsertOne(context.TODO(), s)
	return err
}

func (m *MongoStorage) GetShortURLByID(id string) (*model.ShortURL, error) {
	var s model.ShortURL
	err := m.db.Collection("shorturls").FindOne(context.TODO(), bson.M{"id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *MongoStorage) UpdateShortURL(s model.ShortURL) error {
	_, err := m.db.Collection("shorturls").ReplaceOne(context.TODO(), bson.M{"id": s.ID}, s)
	return err
}

func (m *MongoStorage) DeleteShortURL(id string) error {
	_, err := m.db.Collection("shorturls").DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func (m *MongoStorage) GetAllShortURLs() ([]model.ShortURL, error) {
	cursor, err := m.db.Collection("shorturls").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var shortURLs []model.ShortURL
	for cursor.Next(context.TODO()) {
		var s model.ShortURL
		if err := cursor.Decode(&s); err != nil {
			log.Println("decode error:", err)
			continue
		}
		shortURLs = append(shortURLs, s)
	}
	return shortURLs, nil
}

// Stats

func (m *MongoStorage) SaveStats(s model.Stats) error {
	_, err := m.db.Collection("stats").InsertOne(context.TODO(), s)
	return err
}

func (m *MongoStorage) GetStatsByID(id string) (*model.Stats, error) {
	var s model.Stats
	err := m.db.Collection("stats").FindOne(context.TODO(), bson.M{"id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *MongoStorage) UpdateStats(s model.Stats) error {
	_, err := m.db.Collection("stats").ReplaceOne(context.TODO(), bson.M{"id": s.ID}, s)
	return err
}

func (m *MongoStorage) DeleteStats(id string) error {
	_, err := m.db.Collection("stats").DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func (m *MongoStorage) GetAllStats() ([]model.Stats, error) {
	cursor, err := m.db.Collection("stats").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var stats []model.Stats
	for cursor.Next(context.TODO()) {
		var s model.Stats
		if err := cursor.Decode(&s); err != nil {
			log.Println("decode error:", err)
			continue
		}
		stats = append(stats, s)
	}
	return stats, nil
}
