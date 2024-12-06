package repositories

import (
	"context"
	"errors"

	"ObserverKVS/infrastructure"
	"ObserverKVS/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	CreateUser(user *models.User) (string, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}

type MongoUserContext struct {
	collection *mongo.Collection
}

func NewUserRepository() (*MongoUserContext, error) {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://db:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Get the user collection
	collection := client.Database("observerKVS").Collection("users")

	result, err := insertSeedUser(collection)
	if err != nil || !result {
		return nil, err
	}

	return &MongoUserContext{
		collection: collection,
	}, nil
}

func insertSeedUser(collection *mongo.Collection) (bool, error) {
	// Check if seed data already exists
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return false, err
	}

	// If seed data not exists, create seed data
	if count == 0 {
		apikey, err := infrastructure.Encrypt("kvs_1qaz2wsx3edc4rfv")
		if err != nil {
			return false, err
		}

		// Create seed data
		user := models.User{
			Username: "observer",
			ApiKey:   apikey,
			Role:     "admin",
		}

		_, err = collection.InsertOne(context.Background(), user)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *MongoUserContext) CreateUser(user *models.User) (string, error) {
	// Check if the key already exists in the collection
	existingKeyFilter := bson.M{"username": user.Username}
	existingKeyCount, err := r.collection.CountDocuments(context.Background(), existingKeyFilter)
	if err != nil {
		return "", err
	}

	if existingKeyCount > 0 {
		return "", errors.New("username already exists")
	}

	//Create api key for created user
	user.ApiKey, err = infrastructure.GenerateAPIKey()
	if err != nil {
		return "", err
	}

	user.ApiKey, err = infrastructure.Encrypt(user.ApiKey)
	if err != nil {
		return "", err
	}

	insertResult, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoUserContext) GetUserByID(id string) (*models.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}

	var user models.User
	err = r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, errors.New("notfound")
	}

	user.ApiKey, err = infrastructure.Decrypt(user.ApiKey)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserContext) GetUserByApiKey(apiKey string) (*models.User, error) {
	apiKey, err := infrastructure.Encrypt(apiKey)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"apikey": apiKey}

	var user models.User
	err = r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, errors.New("notfound")
	}

	user.ApiKey, err = infrastructure.Decrypt(user.ApiKey)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserContext) DeleteUser(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectId}

	_, err = r.collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}
	return nil
}
