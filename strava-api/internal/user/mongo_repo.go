package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yellow-jersey/internal/strava"
	"yellow-jersey/pkg/logs"
)

// Ensure MongoRepository implements the Repository interface
var _ Repository = (*MongoRepository)(nil)

// MongoRepository is a Mongo implementation of the user.Repository interface.
type MongoRepository struct {
	db    *mongo.Database
	users *mongo.Collection
}

// NewMongoRepository returns a MongoRepository struct with a database connection.
func NewMongoRepository(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// TODO: Should these come from config
	db := client.Database("yellow-jersey")
	users := db.Collection("users")

	return &MongoRepository{
		db:    db,
		users: users,
	}, nil
}

// NewMongoRepoWithDB returns a new Mongo repo with the supplied Mongo database.
func NewMongoRepoWithDB(db *mongo.Database) (*MongoRepository, error) {
	users := db.Collection("users")
	return &MongoRepository{
		db:    db,
		users: users,
	}, nil
}

// CreateUser creates a new user in the MongoDB.
func (m MongoRepository) CreateUser(accessToken, refreshToken, stravaID string, ath strava.AthleteDetailed) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := NewUser(accessToken, refreshToken, stravaID, ath)
	_, err := m.users.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("problem creating user with Strava ID %s within Mongo : %w", stravaID, err)
	}

	return user, nil
}

// FetchUser fetches a user from Mongo by their internal Yellow Jersey ID.
func (m MongoRepository) FetchUser(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var u *User
	err := m.users.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("unable to fetch user %s from Mongo : %w", id, err)
	}

	return u, nil
}

// FetchUserByStravaID fetches a user by their Strava ID. This method is used to check if the user exists
// when we only have access to their Strava ID.
func (m MongoRepository) FetchUserByStravaID(stravaID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var u *User
	err := m.users.FindOne(ctx, bson.M{"strava_id": stravaID}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("unable to fetch user with Strava ID %s from Mongo : %w", stravaID, err)
	}

	return u, nil
}

// FetchAll fetches all users from the system.
// TODO: Should be limited with pagination in future.
func (m MongoRepository) FetchAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := m.users.Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logs.Logger.Info().Msg("no users found in the db")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to fetch users from Mongo : %w", err)
	}

	var users []*User
	for cur.Next(context.Background()) {
		var user *User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Close(context.Background()); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates the user within Mongo.
func (m MongoRepository) UpdateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.M{
		"$set": user,
	}
	_, err := m.users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
