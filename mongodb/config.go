package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type config struct {
	client *mongo.Client
}
