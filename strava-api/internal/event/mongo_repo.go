package event

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure MongoRepository implements the Repo interface
var _ Repo = (*MongoRepository)(nil)

// MongoRepository is a Mongo implementation of the user.Repository interface.
type MongoRepository struct {
	db     *mongo.Database
	events *mongo.Collection
}

// NewMongoRepository returns a MongoRepository struct with a database connection.
func NewMongoRepository(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// TODO: Should these come from config
	db := client.Database("yellow-jersey")
	users := db.Collection("events")

	return &MongoRepository{
		db:     db,
		events: users,
	}, nil
}

// NewMongoRepoWithDB returns a new Mongo repo with the supplied Mongo database.
func NewMongoRepoWithDB(db *mongo.Database) (*MongoRepository, error) {
	users := db.Collection("events")
	return &MongoRepository{
		db:     db,
		events: users,
	}, nil
}

// Create creates a new event in the MongoDB.
func (m MongoRepository) Create(owner, name string, startDate, finishDate time.Time) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event, err := NewEvent(owner, name, startDate, finishDate)
	if err != nil {
		return nil, err
	}

	_, err = m.events.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("problem creating event for user %s within Mongo : %w", owner, err)
	}

	return event, nil
}

// Fetch fetches an event from Mongo.
func (m MongoRepository) Fetch(id string) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var e *Event
	err := m.events.FindOne(ctx, bson.M{"_id": id}).Decode(&e)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch event %s from Mongo : %w", id, err)
	}

	return e, nil
}

// FetchUserEvents fetches all events for a user.
func (m MongoRepository) FetchUserEvents(userID string) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := m.events.Find(ctx, bson.M{"owner": userID})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch events for user %s from Mongo : %w", userID, err)
	}
	defer cur.Close(ctx)

	var events []*Event
	err = cur.All(ctx, &events)

	return events, err
}

// Update updates an event within Mongo.
func (m MongoRepository) Update(e *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: e.ID}}
	update := bson.M{
		"$set": e,
	}
	_, err := m.events.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
