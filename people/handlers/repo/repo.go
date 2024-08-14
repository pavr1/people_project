package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/pavr1/people_project/people/config"
	"github.com/pavr1/people_project/people/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepoHandler struct {
	log    *log.Logger
	Config *config.Config
	client *mongo.Client
}

func NewRepoHandler(log *log.Logger, config *config.Config) (*RepoHandler, error) {
	client, err := connectToMongoDB(config)
	if err != nil {
		log.WithField("error", err).Error("Failed to connect to MongoDB")

		return nil, err
	}

	return &RepoHandler{
		log:    log,
		Config: config,
		client: client,
	}, nil
}

func connectToMongoDB(config *config.Config) (*mongo.Client, error) {
	uri := config.MongoDB.Uri

	log.WithField("uri", uri).Info("Connecting to MongoDB...")

	clientOptions := options.Client().ApplyURI(uri)
	// clientOptions.SetAuth(options.Credential{
	// 	Username: config.MongoDB.Username,
	// 	Password: config.MongoDB.Password,
	// })
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.WithError(err).Error("Failed to connect to MongoDB")

		return nil, err
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func (r *RepoHandler) GetPersonList() ([]models.Person, error) {
	people := []models.Person{}

	// Get a handle to the collection
	collection := r.client.Database(r.Config.MongoDB.Database).Collection(r.Config.MongoDB.Collection)

	// Find all documents in the collection
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.WithError(err).Error("Failed to find documents in MongoDB")

		return nil, err
	}

	defer cur.Close(context.Background())

	// Iterate over the documents and print their contents
	for cur.Next(context.Background()) {
		var doc bson.M
		err := cur.Decode(&doc)
		if err != nil {
			log.Error(err)
		}

		people = append(people, models.Person{
			ID:       doc["id"].(string),
			Name:     doc["name"].(string),
			LastName: doc["lastName"].(string),
			Age:      doc["age"].(int32),
		})
	}

	if err := cur.Err(); err != nil {
		log.WithError(err).Error("Failed to iterate over documents in MongoDB")

		return nil, err
	}

	return people, nil
}

func (r *RepoHandler) GetPerson(id string) (*models.Person, error) {
	// Get the database and collection
	db := r.client.Database(r.Config.MongoDB.Database)
	collection := db.Collection(r.Config.MongoDB.Collection)

	// Find the document by ID
	filter := bson.M{"id": id}
	var person models.Person
	err := collection.FindOne(context.Background(), filter).Decode(&person)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		log.WithError(err).Error("Failed to find document in MongoDB")

		return nil, err
	}

	return &person, nil
}

func (r *RepoHandler) CreatePerson(person *models.Person) error {
	existentPerson, err := r.GetPerson(person.ID)
	if err != nil {
		//will need to check for not found
		log.WithError(err).Error("Failed to get person from MongoDB")

		return err
	}

	if existentPerson != nil {
		log.WithField("id", person.ID).Info("Person already exists")

		return fmt.Errorf("person with ID %s already exists", person.ID)
	}

	// Insert the person into the "people" collection
	collection := r.client.Database(r.Config.MongoDB.Database).Collection(r.Config.MongoDB.Collection)

	doc := bson.D{}

	// Add fields to the document
	doc = append(doc, bson.E{Key: "id", Value: person.ID})
	doc = append(doc, bson.E{Key: "name", Value: person.Name})
	doc = append(doc, bson.E{Key: "lastName", Value: person.LastName})
	doc = append(doc, bson.E{Key: "age", Value: person.Age})

	// Convert the document to BSON
	personBSON, err := bson.Marshal(doc)
	if err != nil {
		log.WithError(err).Error("Failed to marshal person to BSON")
		return err
	}

	_, err = collection.InsertOne(context.Background(), personBSON)
	if err != nil {
		log.WithError(err).Error("Failed to insert person into MongoDB")

		return err
	}

	log.WithField("id", person.ID).Info("Person inserted successfully")

	return nil
}

func (r *RepoHandler) DeletePerson(id string) error {
	person, err := r.GetPerson(id)
	if err != nil {
		return err
	}

	if person == nil {
		return fmt.Errorf("person with ID %s not found", id)
	}

	// Get the database and collection
	db := r.client.Database(r.Config.MongoDB.Database)
	collection := db.Collection(r.Config.MongoDB.Collection)

	// Delete the document by ID
	filter := bson.M{"id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.WithError(err).Error("Failed to delete document from MongoDB")

		return err
	}

	log.WithField("id", id).Info("Person deleted successfully")

	return nil
}

func (r *RepoHandler) UpdatePerson(person *models.Person) error {
	// Get the database and collection
	db := r.client.Database(r.Config.MongoDB.Database)
	collection := db.Collection(r.Config.MongoDB.Collection)

	// Update the document by ID
	filter := bson.M{"id": person.ID}
	update := bson.M{"$set": person}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.WithError(err).Error("Failed to update document in MongoDB")

		return err
	}

	log.WithField("id", person.ID).Info("Person updated successfully")

	return nil
}
