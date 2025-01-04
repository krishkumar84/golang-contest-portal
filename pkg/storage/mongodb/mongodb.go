package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
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

func (m *MongoDB) GetAllContests() ([]types.ContestBasicInfo, error) {
    collection := m.db.Collection("contests")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    projection := bson.D{
        {Key: "_id", Value: 1},
        {Key: "title", Value: 1},
        {Key: "start_time", Value: 1},
        {Key: "end_time", Value: 1},
        {Key: "description", Value: 1},
    }

    var contests []types.ContestBasicInfo
    cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &contests); err != nil {
        return nil, err
    }

    return contests, nil
}

func (m *MongoDB) GetContestById(id string) ([]bson.M, error) {
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("invalid contest id format")
    }

    pipeline := mongo.Pipeline{
        {{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},
        {{Key: "$addFields", Value: bson.D{
            {Key: "question_ids", Value: bson.D{
                {Key: "$map", Value: bson.D{
                    {Key: "input", Value: "$question_ids"},
                    {Key: "as", Value: "qid"},
                    {Key: "in", Value: bson.D{
                        {Key: "$toObjectId", Value: "$$qid"},
                    }},
                }},
            }},
        }}},
        {{Key: "$lookup", Value: bson.D{
            {Key: "from", Value: "questions"},
            {Key: "localField", Value: "question_ids"},
            {Key: "foreignField", Value: "_id"},
            {Key: "as", Value: "questions"},
        }}},
        {{Key: "$project", Value: bson.D{
            {Key: "_id", Value: "$_id"},
            {Key: "title", Value: "$title"},
            {Key: "start_time", Value: "$start_time"},
            {Key: "end_time", Value: "$end_time"},
            {Key: "description", Value: "$description"},
            {Key: "questions", Value: bson.D{
                {Key: "$map", Value: bson.D{
                    {Key: "input", Value: "$questions"},
                    {Key: "as", Value: "q"},
                    {Key: "in", Value: bson.D{
                        {Key: "_id", Value: "$$q._id"},
                        {Key: "title", Value: "$$q.title"},
                        {Key: "description", Value: "$$q.description"},
                        {Key: "difficulty", Value: "$$q.difficulty"},
                    }},
                }},
            }},
        }}},
    }

    var results []bson.M
    cursor, err := m.db.Collection("contests").Aggregate(context.Background(), pipeline)
    if err != nil {
        return nil, fmt.Errorf("error executing aggregation: %v", err)
    }
    defer cursor.Close(context.Background())

    if err := cursor.All(context.Background(), &results); err != nil {
        return nil, fmt.Errorf("error decoding result: %v", err)
    }

    if len(results) == 0 {
        return nil, mongo.ErrNoDocuments
    }

    return results, nil
}