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

type KeyValuePairRepository interface {
	Save(pair *models.KeyValuePair) (string, error)
	Get(key string) (string, error)
	GetById(id string) (string, error)
	Update(pair *models.KeyValuePair) (string, error)
	Delete(key string) error
}

type MongoKeyValuePairContext struct {
	collection *mongo.Collection
}

func NewKeyValuePairRepository(stackholder string) (*MongoKeyValuePairContext, error) {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Get the key-value pair collection
	collection := client.Database("observerKVS").Collection("keyvaluepairs_" + stackholder)

	return &MongoKeyValuePairContext{
		collection: collection,
	}, nil
}

func (r *MongoKeyValuePairContext) Save(pair *models.KeyValuePair) (string, error) {
	// Create a key-value pair document
	// Check if the key already exists in the collection
	existingKeyFilter := bson.M{"key": pair.Key}
	existingKeyCount, err := r.collection.CountDocuments(context.Background(), existingKeyFilter)
	if err != nil {
		return "", err
	}
	if existingKeyCount > 0 {
		return "", errors.New("key already exists")
	}

	value, err := infrastructure.Encrypt(pair.Value)
	if err != nil {
		return "", err
	}

	document := bson.M{
		"key":   pair.Key,
		"value": value,
	}

	// Insert the document into the collection
	insertResult, err := r.collection.InsertOne(context.Background(), document)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoKeyValuePairContext) Get(key string) (string, error) {
	// Define a filter for the key
	filter := bson.M{"key": key}

	// Find the document with the specified key
	var result bson.M
	err := r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", errors.New("notfound")
	}

	// Extract the value from the document
	value := result["value"].(string)

	value, err = infrastructure.Decrypt(value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *MongoKeyValuePairContext) GetById(id string) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": objectId}

	// Find the document with the specified key
	var result bson.M
	err = r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", errors.New("notfound")
	}

	// Extract the value from the document
	value := result["value"].(string)

	value, err = infrastructure.Decrypt(value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *MongoKeyValuePairContext) Update(pair *models.KeyValuePair) (bool, error) {
	// Define a filter for the key
	filter := bson.M{"key": pair.Key}

	value, err := infrastructure.Encrypt(pair.Value)
	if err != nil {
		return false, err
	}

	// Define an update document
	update := bson.M{"$set": bson.M{"value": value}}

	// Update the document in the collection
	updateResult, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, err
	}

	return updateResult.ModifiedCount > 0, nil
}

func (r *MongoKeyValuePairContext) Delete(key string) error {
	// Define a filter for the key
	filter := bson.M{"key": key}

	// Delete the document from the collection
	_, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
