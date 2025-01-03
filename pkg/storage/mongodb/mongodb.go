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
    "github.com/go-playground/validator/v10"
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
func (m *MongoDB) CreateUser(name, email, password, studentId string, role types.Role) (string, error) {
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
        Role:      role,
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

func (m *MongoDB) CreateContest(contest types.Contest) (string, error) {
    collection := m.db.Collection("contests")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := validator.New().Struct(contest); err != nil {
        validateErrs := err.(validator.ValidationErrors)
        return "", fmt.Errorf("validation failed: %v", validateErrs)
    }
    result, err := collection.InsertOne(ctx, contest)
    if err != nil {
        return "", err
    }

    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *MongoDB) CreateQuestion(question types.Question) (string, error) {
    collection := m.db.Collection("questions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := validator.New().Struct(question); err != nil {
        validateErrs := err.(validator.ValidationErrors)
        return "", fmt.Errorf("validation failed: %v", validateErrs)
    }

    result, err := collection.InsertOne(ctx, question)
    if err != nil {
        return "", err
    }

    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *MongoDB) CreateTestCase(testCase types.TestCase) (string, error) {
    collection := m.db.Collection("test_cases")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := validator.New().Struct(testCase); err != nil {
        validateErrs := err.(validator.ValidationErrors)
        return "", fmt.Errorf("validation failed: %v", validateErrs)
    }

    result, err := collection.InsertOne(ctx, testCase)
    if err != nil {
        return "", err
    }

    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}