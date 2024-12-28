package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/config"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	db *mongo.Database
}

func New(cfg *config.Config) (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(cfg.DatabaseURL)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }

    db := client.Database(cfg.DatabaseName)
    return &MongoDB{
        client: client,
        db:     db,
    }, nil
}

// User operations
func (m *MongoDB) CreateUser(name, email, password, studentId string) (string, error) {
    collection := m.db.Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    emailCount, err := collection.CountDocuments(ctx, bson.M{"email": email})
    if err != nil {
        return "", err
    }
    if emailCount > 0 {
        return "", fmt.Errorf("user with email %s already exists", email)
    }

    studentIdCount, err := collection.CountDocuments(ctx, bson.M{"studentId": studentId})
    if err != nil {
        return "", err
    }
    if studentIdCount > 0 {
        return "", fmt.Errorf("user with student ID %s already exists", studentId)
    }

    user := types.User{
        ID:        primitive.NewObjectID(),
        Name:      name,
        Email:     email,
        StudentId: studentId,
        Password:  password,
        CreatedAt: time.Now(),
    }

    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        return "", err
    }

    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}



func (m *MongoDB) GetUserByEmail(email string) (*types.User, error) {
    collection := m.db.Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user types.User
    err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}