package event_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"yellow-jersey/internal/event"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

const MongoInitdbRootUsername = "root"
const MongoInitdbRootPassword = "password"

func TestMain(m *testing.M) {
	var client *mongo.Client
	// Setup
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	environmentVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=" + MongoInitdbRootUsername,
		"MONGO_INITDB_ROOT_PASSWORD=" + MongoInitdbRootPassword,
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env:        environmentVariables,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		client, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%s@localhost:%s", MongoInitdbRootUsername, MongoInitdbRootPassword, resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		db = client.Database("yellow-jersey")
		return client.Ping(context.TODO(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run tests
	exitCode := m.Run()

	// Teardown
	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	// Exit
	os.Exit(exitCode)
}

func TestMongoRepository_Create(t *testing.T) {
	mongo, err := event.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	event, err := mongo.Create("test", "event_name", time.Now(), time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, "test", event.Owner)
	assert.Equal(t, "event_name", event.Name)
}

func TestMongoRepository_Fetch(t *testing.T) {
	mongo, err := event.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	event, err := mongo.Create("test", "event_name", time.Now(), time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, "test", event.Owner)

	returnedEvent, err := mongo.Fetch(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, event.ID, returnedEvent.ID)
}

func TestMongoRepository_Update(t *testing.T) {
	mongo, err := event.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	event, err := mongo.Create("test", "event_name", time.Now(), time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, "test", event.Owner)

	event.Owner = "new"
	assert.NoError(t, mongo.Update(event))

	returnedEvent, err := mongo.Fetch(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, event.Owner, returnedEvent.Owner)
}
