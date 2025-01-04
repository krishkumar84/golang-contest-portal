package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Role string
const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name      string            `bson:"name" json:"name" validate:"required"`
    Email     string            `bson:"email" json:"email" validate:"required"`
    Password  string            `bson:"password" json:"password" validate:"required"`
    StudentId string            `bson:"studentId" json:"studentId" validate:"required"`
    CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
    Role      Role              `bson:"role" json:"role" default:"user"`
}

type Contest struct {
    ID          primitive.ObjectID   `bson:"_id,omitempty" json:"contest_id"`
    Title       string              `bson:"title" json:"title" validate:"required"`
    StartTime   time.Time           `bson:"start_time" json:"start_time" validate:"required"`
    EndTime     time.Time           `bson:"end_time" json:"end_time" validate:"required"`
    Description string              `bson:"description" json:"description" validate:"required"`
    CreatedBy   string              `bson:"created_by" json:"created_by"`
    QuestionIDs []string            `bson:"question_ids" json:"question_ids,omitempty"`
    CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
}

type Question struct {
    ID        string    `bson:"_id,omitempty" json:"question_id"`
    Title     string    `bson:"title" json:"title" validate:"required"`
    Description string `bson:"description" json:"description" validate:"required"`
    Difficulty string `bson:"difficulty" json:"difficulty"`
    Tags []string     `bson:"tags" json:"tags"`
    TestCaseIDs []string `bson:"test_case_ids" json:"test_case_ids" validate:"required"`
    Points int `bson:"points" json:"points"`
    Cpu_time_limit int `bson:"cpu_time_limit" json:"cpu_time_limit"`
    Memory_limit int `bson:"memory_limit" json:"memory_limit"`
    CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Visibility string

const (
    VisibilityPublic  Visibility = "public"
    VisibilityPrivate Visibility = "private"
)

type TestCase struct {
    ID string `bson:"_id,omitempty" json:"test_case_id"`
    Input interface{} `bson:"input" json:"input"`
    ExpectedOutput interface{} `bson:"expected_output" json:"expected_output"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    Visibility Visibility `bson:"visibility" json:"visibility" validate:"required,oneof=public private"`
}

type Submission struct {
    ID string `bson:"_id,omitempty" json:"submission_id"`
    UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
    QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
    ContestID primitive.ObjectID `bson:"contest_id" json:"contest_id"`
    Code string `bson:"code" json:"code"`
    LanguageId string `bson:"language_id" json:"language_id"`
    Status string `bson:"status" json:"status"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Leaderboard struct {
    ID string `bson:"_id,omitempty" json:"leaderboard_id"`
    UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
    ContestID primitive.ObjectID `bson:"contest_id" json:"contest_id"`
    LeaderboardScore int `bson:"leaderboard_score" json:"leaderboard_score"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}