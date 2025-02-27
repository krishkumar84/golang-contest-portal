package storage

import (
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
)


type Storage interface {
	// CreateUser(name string, email string, age int)(int64, error)
	// GetUserById(id int64)(types.User, error) 
	// GetAllUsers()([]types.User, error)
	CreateUser(name string, email string, password string, studentId string, role types.Role) (string, error)
	GetUserByEmail(email string) (*types.User, error)
	CreateContest(contest types.Contest) (string, error)
	DeleteContestById(id string) error
	CreateQuestion(question types.Question) (string, error)
	EditQuestionById(id string, question types.Question) error
	EditContestById(id string, contest types.Contest) error
	DeleteQuestionFromContestById(contestId string, questionId string) error
	CreateTestCase(testCase types.TestCase) (string, error)
	GetAllContests() ([]types.ContestBasicInfo, error)
	GetContestById(id string) ([]bson.M, error)
	GetQuestionById(id string) ([]bson.M, error)
	AddQuestionToContest(contestId string, question types.Question) (string, error)
	AddTestCaseToQuestion(questionId string, testCase types.TestCase) (string, error)
	DeleteTestCaseFromQuestionById(questionId string, testCaseId string) error
	EditTestCaseById(testCaseId string, testCase types.TestCase) error
	CreateSubmission(submission types.Submission) (string, error)
	GetSubmissionById(id string) (*types.Submission, error)
	UpdateSubmissionStatus(id string, status string, score int) error
}
