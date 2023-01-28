package user_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yellow-jersey/internal/user"
)

var db *mongo.Database

const MongoInitdbRootUsername = "root"
const MongoInitdbRootPassword = "password"

// TODO: Move all this into a central place for use in multiple packages
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
				fmt.Sprintf(
					"mongodb://%s:%s@localhost:%s",
					MongoInitdbRootUsername,
					MongoInitdbRootPassword,
					resource.GetPort("27017/tcp"),
				),
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
	mongo, err := user.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	usr, err := mongo.CreateUser("token", "refresh_token", "strava_id")
	assert.NoError(t, err)
	assert.Equal(t, "strava_id", usr.StravaID)
	assert.Equal(t, "token", usr.AccessToken)
	assert.Equal(t, "refresh_token", usr.RefreshToken)
}

func TestMongoRepository_Fetch(t *testing.T) {
	mongo, err := user.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	usr, err := mongo.CreateUser("token", "refresh_token", "strava_id")
	assert.NoError(t, err)

	_, err = mongo.FetchUser(usr.ID)
	assert.NoError(t, err)
}

func TestMongoRepository_Fetch_ByStravaID(t *testing.T) {
	mongo, err := user.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	usr, err := mongo.CreateUser("token", "refresh_token", "strava_id")
	assert.NoError(t, err)

	usr, err = mongo.FetchUserByStravaID(usr.StravaID)
	assert.NoError(t, err)
	assert.Equal(t, "refresh_token", usr.RefreshToken)
}

func TestMongoRepository_Update(t *testing.T) {
	mongo, err := user.NewMongoRepoWithDB(db)
	assert.NoError(t, err)

	usr, err := mongo.CreateUser("token", "refresh_token", "strava_id")
	assert.NoError(t, err)

	usr.AccessToken = "new"
	assert.NoError(t, mongo.UpdateUser(usr))

	returnedUsr, err := mongo.FetchUser(usr.ID)
	assert.NoError(t, err)
	assert.Equal(t, returnedUsr.AccessToken, "new")
}
