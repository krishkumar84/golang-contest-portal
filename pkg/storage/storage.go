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
	CreateQuestion(question types.Question) (string, error)
	CreateTestCase(testCase types.TestCase) (string, error)
	GetAllContests() ([]types.ContestBasicInfo, error)
	GetContestById(id string) ([]bson.M, error)
}
