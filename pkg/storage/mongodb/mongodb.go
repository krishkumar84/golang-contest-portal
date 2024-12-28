package mongodb

import (
    "context"
    "time"
    
    "github.com/krishkumar84/bdcoe-golang-portal/pkg/config"
    "github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
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
func (m *MongoDB) CreateUser(name, email, password string) (string, error) {
    collection := m.db.Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    user := types.User{
        Name:      name,
        Email:     email,
        Password:  password,  // Note: In production, ensure this is hashed
        CreatedAt: time.Now(),
    }

    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        return "", err
    }

    return result.InsertedID.(string), nil
}