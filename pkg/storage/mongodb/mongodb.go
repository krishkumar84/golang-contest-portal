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

    if contest.QuestionIDs == nil {
        contest.QuestionIDs = []string{}
    }

    result, err := collection.InsertOne(ctx, contest)
    if err != nil {
        return "", err
    }

    return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *MongoDB) EditContestById(id string, updateData types.Contest) error {
    contestObjID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid contest id format")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var contest types.Contest
    err = m.db.Collection("contests").FindOne(ctx, bson.M{"_id": contestObjID}).Decode(&contest)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fmt.Errorf("no contest found with the given id")
        }
        return fmt.Errorf("error checking contest existence: %v", err)
    }

    // Prepare the update
    update := bson.M{}
    if updateData.Title != "" {
        update["title"] = updateData.Title
    }
    if !updateData.StartTime.IsZero() {
        update["start_time"] = updateData.StartTime
    }
    if !updateData.EndTime.IsZero() {
        update["end_time"] = updateData.EndTime
    }
    if updateData.Description != "" {
        update["description"] = updateData.Description
    }
    if updateData.CreatedBy != "" {
        update["created_by"] = updateData.CreatedBy
    }

    if len(update) > 0 {
        _, err = m.db.Collection("contests").UpdateOne(ctx, bson.M{"_id": contestObjID}, bson.M{"$set": update})
        if err != nil {
            return fmt.Errorf("failed to update contest: %v", err)
        }
    }

    return nil
}

func (m *MongoDB) DeleteContestById(id string) error {
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid contest id format")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := m.db.Collection("contests").DeleteOne(ctx, bson.M{"_id": objectId})
    if err != nil {
        return err
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("no contest found with the given id")
    }

    return nil
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

func (m *MongoDB) GetQuestionById(id string) ([]bson.M, error) {
    collection := m.db.Collection("questions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("invalid question id format")
    }

    pipeline := mongo.Pipeline{
        {{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},
        {{Key: "$addFields", Value: bson.D{
            {Key: "test_case_ids", Value: bson.D{
                {Key: "$map", Value: bson.D{
                    {Key: "input", Value: "$test_case_ids"},
                    {Key: "as", Value: "tid"},
                    {Key: "in", Value: bson.D{
                        {Key: "$toObjectId", Value: "$$tid"},
                    }},
                }},
            }},
        }}},
        {{Key: "$lookup", Value: bson.D{
            {Key: "from", Value: "test_cases"},
            {Key: "localField", Value: "test_case_ids"},
            {Key: "foreignField", Value: "_id"},
            {Key: "as", Value: "test_cases"},
        }}},
        {{Key: "$project", Value: bson.D{
            {Key: "_id", Value: 1},
            {Key: "title", Value: 1},
            {Key: "description", Value: 1},
            {Key: "difficulty", Value: 1},
            {Key: "tags", Value: 1},
            {Key: "points", Value: 1},
            {Key: "test_cases", Value: bson.D{
                {Key: "$map", Value: bson.D{
                    {Key: "input", Value: "$test_cases"},
                    {Key: "as", Value: "tc"},
                    {Key: "in", Value: bson.D{
                        {Key: "_id", Value: "$$tc._id"},
                        {Key: "input", Value: "$$tc.input"},
                        {Key: "expected_output", Value: "$$tc.expected_output"},
                        {Key: "visibility", Value: "$$tc.visibility"},
                    }},
                }},
            }},
        }}},
    }

    var results []bson.M
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("error executing aggregation: %v", err)
    }
    defer cursor.Close(ctx)

    if err := cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("error decoding result: %v", err)
    }

    if len(results) == 0 {
        return nil, mongo.ErrNoDocuments
    }

    return results, nil
}

func (m*MongoDB)  EditQuestionById(id string, updateData types.Question) error {
    questionObjID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid question id format")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var question types.Question
    err = m.db.Collection("questions").FindOne(ctx, bson.M{"_id": questionObjID}).Decode(&question)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fmt.Errorf("no question found with the given id")
        }
        return fmt.Errorf("error checking question existence: %v", err)
    }

    // Prepare the update
    update := bson.M{}
    if updateData.Title != "" {
        update["title"] = updateData.Title
    }
    if updateData.Description != "" {
        update["description"] = updateData.Description
    }
    if updateData.Difficulty != "" {
        update["difficulty"] = updateData.Difficulty
    }
    if updateData.Tags != nil {
        update["tags"] = updateData.Tags
    }
    if updateData.Points != 0 {
        update["points"] = updateData.Points
    }
    if updateData.Cpu_time_limit != 0 {
        update["cpu_time_limit"] = updateData.Cpu_time_limit
    }
    if updateData.Memory_limit != 0 {
        update["memory_limit"] = updateData.Memory_limit
    }
    if len(update) > 0 {
        _, err = m.db.Collection("questions").UpdateOne(ctx, bson.M{"_id": questionObjID}, bson.M{"$set": update})
        if err != nil {
            return fmt.Errorf("failed to update question: %v", err)
        }
    }

    return nil
}

func (m *MongoDB) AddQuestionToContest(contestId string, question types.Question) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    fmt.Printf("Received contest ID: %s\n", contestId)

    question.TestCaseIDs = []string{}

    //Create question
    questionResult, err := m.db.Collection("questions").InsertOne(ctx, question)
    if err != nil {
        return "", fmt.Errorf("failed to create question: %v", err)
    }
    questionId := questionResult.InsertedID.(primitive.ObjectID).Hex()

    fmt.Printf("Created question with ID: %s\n", questionId)

    contestObjID, err := primitive.ObjectIDFromHex(contestId)
    if err != nil {
        return "", fmt.Errorf("invalid contest id format: %v", err)
    }

    filter := bson.M{"_id": contestObjID}
    update := bson.M{"$push": bson.M{"question_ids": questionId}}
    
    result, err := m.db.Collection("contests").UpdateOne(ctx, filter, update)
    if err != nil {
        return "", fmt.Errorf("failed to update contest: %v", err)
    }

    fmt.Printf("Updated contest. Modified count: %d\n", result.ModifiedCount)

    return questionId, nil
}

func (m *MongoDB) DeleteQuestionFromContestById(contestId string, questionId string) error {
    contestObjID, err := primitive.ObjectIDFromHex(contestId)
    if err != nil {
        return fmt.Errorf("invalid contest id format")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if contest exists
    var contest types.Contest
    err = m.db.Collection("contests").FindOne(ctx, bson.M{"_id": contestObjID}).Decode(&contest)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fmt.Errorf("no contest found with the given id")
        }
        return fmt.Errorf("error checking contest existence: %v", err)
    }

    found := false
    for _, qID := range contest.QuestionIDs {
        if qID == questionId {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("no question found with the given id in the contest")
    }

    filter := bson.M{"_id": contestObjID}
    update := bson.M{"$pull": bson.M{"question_ids": questionId}}
    
    result, err := m.db.Collection("contests").UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update contest: %v", err)
    }

    if result.ModifiedCount == 0 {
        return fmt.Errorf("no question found with the given id in the contest")
    }

    return nil
}

func (m *MongoDB) AddTestCaseToQuestion(questionId string, testCase types.TestCase) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    fmt.Printf("Received question ID: %s\n", questionId)

    // Create test case
    testCaseResult, err := m.db.Collection("test_cases").InsertOne(ctx, testCase)
    if err != nil {
        return "", fmt.Errorf("failed to create test case: %v", err)
    }
    testCaseId := testCaseResult.InsertedID.(primitive.ObjectID).Hex()

    fmt.Printf("Created test case with ID: %s\n", testCaseId)

    // Convert question ID to ObjectID
    questionObjID, err := primitive.ObjectIDFromHex(questionId)
    if err != nil {
        return "", fmt.Errorf("invalid question id format: %v", err)
    }

    // Update question with test case ID
    filter := bson.M{"_id": questionObjID}
    update := bson.M{"$push": bson.M{"test_case_ids": testCaseId}}
    
    result, err := m.db.Collection("questions").UpdateOne(ctx, filter, update)
    if err != nil {
        return "", fmt.Errorf("failed to update question: %v", err)
    }

    fmt.Printf("Updated question. Modified count: %d\n", result.ModifiedCount)

    return testCaseId, nil
}

func (m *MongoDB) DeleteTestCaseFromQuestionById(questionId string, testCaseId string) error {
    questionObjID, err := primitive.ObjectIDFromHex(questionId)
    if err != nil {
        return fmt.Errorf("invalid question id format")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if question exists
    var question types.Question
    err = m.db.Collection("questions").FindOne(ctx, bson.M{"_id": questionObjID}).Decode(&question)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fmt.Errorf("no question found with the given id")
        }
        return fmt.Errorf("error checking question existence: %v", err)
    }

    found := false
    for _, tcID := range question.TestCaseIDs {
        if tcID == testCaseId {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("no test case found with the given id in the question")
    }

    filter := bson.M{"_id": questionObjID}
    update := bson.M{"$pull": bson.M{"test_case_ids": testCaseId}}
    
    result, err := m.db.Collection("questions").UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update question: %v", err)
    }

    if result.ModifiedCount == 0 {
        return fmt.Errorf("no test case found with the given id in the question")
    }

    return nil
}